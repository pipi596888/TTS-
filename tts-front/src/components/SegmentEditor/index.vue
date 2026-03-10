<template>
  <div class="segment-editor">
    <el-table
      v-if="segments.length > 0"
      :data="segments"
      row-key="id"
      stripe
      border
    >
      <el-table-column label="序号" width="60" align="center">
        <template #default="{ $index }">
          <span class="index-num">{{ $index + 1 }}</span>
        </template>
      </el-table-column>
      <el-table-column label="音色" width="150">
        <template #default="{ row }">
          <VoiceSelector
            v-model="row.voiceId"
            @update:modelValue="(v) => handleUpdate(row.id, { voiceId: v })"
          />
        </template>
      </el-table-column>
      <el-table-column label="情绪" width="120">
        <template #default="{ row }">
          <el-select
            v-model="row.emotion"
            placeholder="情绪"
            clearable
            size="small"
            @change="handleUpdate(row.id, { emotion: row.emotion ?? '' })"
          >
            <el-option label="中性" value="neutral" />
            <el-option label="开心" value="happy" />
            <el-option label="悲伤" value="sad" />
            <el-option label="愤怒" value="angry" />
          </el-select>
        </template>
      </el-table-column>
      <el-table-column label="文本内容">
        <template #default="{ row }">
          <el-input
            v-model="row.text"
            type="textarea"
            :rows="2"
            placeholder="请输入要转换的文本..."
            @input="handleUpdate(row.id, { text: row.text })"
          />
        </template>
      </el-table-column>
      <el-table-column label="字符" width="60" align="center">
        <template #default="{ row }">
          <span class="char-count">{{ row.text.length }}</span>
        </template>
      </el-table-column>
      <el-table-column label="操作" width="70" align="center">
        <template #default="{ row }">
          <el-button
            type="danger"
            :icon="Delete"
            circle
            size="small"
            @click="handleDelete(row.id)"
          />
        </template>
      </el-table-column>
    </el-table>
    <div v-else class="empty-tip">
      <el-empty description="暂无文本片段，请点击右上角添加" :image-size="100" />
    </div>
  </div>
</template>

<script setup lang="ts">
import { Delete } from '@element-plus/icons-vue'
import VoiceSelector from '@/components/VoiceSelector/index.vue'
import { useTTSStore } from '@/store/tts'
import { storeToRefs } from 'pinia'

const ttsStore = useTTSStore()
const { segments } = storeToRefs(ttsStore)

function handleUpdate(id: string, updates: any) {
  ttsStore.updateSegment(id, updates)
}

function handleDelete(id: string) {
  ttsStore.removeSegment(id)
}
</script>

<style scoped>
.segment-editor {
  width: 100%;
  height: 100%;
}

.segment-editor :deep(.el-table) {
  font-size: 13px;
}

.segment-editor :deep(.el-textarea__inner) {
  font-size: 13px;
}

.index-num {
  font-weight: 600;
  color: #409eff;
}

.char-count {
  font-weight: 600;
  color: #67c23a;
}

.empty-tip {
  padding: 40px 0;
}
</style>

