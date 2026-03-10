package worker

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"time"

	"tts-backend/tts-worker/internal/config"
	"tts-backend/tts-worker/internal/engine"
	"tts-backend/tts-worker/internal/model"
	"tts-backend/tts-worker/internal/utils"
)

type TTSWorker struct {
	config       *config.Config
	taskModel    model.TtsTaskModel
	segmentModel model.TtsSegmentModel
	engine       engine.TTSProvider
	merger       *utils.AudioMerger
}

type TaskMessage struct {
	TaskId string `json:"taskId"`
}

func NewTTSWorker(c *config.Config) *TTSWorker {
	db, err := sql.Open("mysql", c.Mysql.DataSource)
	if err != nil {
		panic(err)
	}

	var ttsEngine engine.TTSProvider
	if c.Aliyun.AccessKeyId != "" {
		ttsEngine = engine.NewMockEngine()
	} else {
		ttsEngine = engine.NewMockEngine()
	}

	return &TTSWorker{
		config:       c,
		taskModel:    model.NewTtsTaskModel(db),
		segmentModel: model.NewTtsSegmentModel(db),
		engine:       ttsEngine,
		merger:       utils.NewAudioMerger(),
	}
}

func (w *TTSWorker) ProcessTask(taskId string) error {
	log.Printf("Processing task: %s", taskId)

	err := w.taskModel.UpdateStatus(taskId, "processing", 0)
	if err != nil {
		return err
	}

	_, err = w.taskModel.FindByTaskId(taskId)
	if err != nil {
		w.taskModel.UpdateError(taskId, err.Error())
		return err
	}

	segments, err := w.segmentModel.FindByTaskId(taskId)
	if err != nil {
		w.taskModel.UpdateError(taskId, err.Error())
		return err
	}

	if len(segments) == 0 {
		w.taskModel.UpdateError(taskId, "no segments found")
		return fmt.Errorf("no segments found")
	}

	// 使用第一个片段生成音频
	seg := segments[0]
	log.Printf("Generating audio for text: %s", seg.Text)

	// MockEngine 返回音频 URL
	audioData, err := w.engine.Generate(seg.Text, seg.VoiceId, seg.Emotion)
	if err != nil {
		w.taskModel.UpdateError(taskId, err.Error())
		return err
	}

	// 更新进度
	w.taskModel.UpdateStatus(taskId, "processing", 80)

	// 使用引擎返回的 URL
	audioUrl := string(audioData)
	if audioUrl == "" {
		audioUrl = "https://www2.cs.uic.edu/~i101/SoundFiles/BabyElephantWalk60.wav"
	}

	log.Printf("Task %s completed, audio URL: %s", taskId, audioUrl)

	w.taskModel.UpdateAudioUrl(taskId, audioUrl)

	return nil
}

func (w *TTSWorker) Start(ctx context.Context) error {
	log.Println("TTS Worker started")

	ticker := time.NewTicker(5 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			log.Println("Worker stopping")
			return nil
		case <-ticker.C:
			w.processPendingTasks()
		}
	}
}

func (w *TTSWorker) processPendingTasks() {
	tasks, err := w.taskModel.FindPendingTasks(5)
	if err != nil {
		log.Printf("Failed to find pending tasks: %v", err)
		return
	}

	if len(tasks) == 0 {
		return
	}

	log.Printf("Found %d pending tasks", len(tasks))

	for _, task := range tasks {
		err := w.ProcessTask(task.TaskId)
		if err != nil {
			log.Printf("Failed to process task %s: %v", task.TaskId, err)
		}
	}
}

func HandleTaskMessage(data []byte) error {
	var msg TaskMessage
	err := json.Unmarshal(data, &msg)
	if err != nil {
		return err
	}

	fmt.Printf("Received task: %s\n", msg.TaskId)
	return nil
}
