<template>
  <div>
    <v-file-input
      v-model="inputImages"
      label="画像"
      filled
      multiple
      clearable
      chips
      accept="image/*"
      placeholder="画像を添付"
      @update:model-value="imageChange"
    />
    <div v-for="(imageUrl, index) in uploadImageUrl" :key="index">
      <v-img :src="imageUrl" max-width="50%" />
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, nextTick } from "vue";

const images = defineModel<File[]>({ default: () => [] });
const inputImages = ref<File[]>([]);
const uploadImageUrl = ref<string[]>([]);

const imageChange = (files: File | File[]) => {
  if (!files) return;
  const fileList = Array.isArray(files) ? files : [files];
  if (fileList.length === 0) return;

  const newImages = [...images.value];

  fileList.forEach((file: File) => {
    const fr = new FileReader();
    fr.readAsDataURL(file);
    newImages.push(file);
    fr.addEventListener("load", () => {
      uploadImageUrl.value.push(fr.result as string);
    });
  });

  images.value = newImages;

  nextTick(() => {
    inputImages.value = [];
  });
};
</script>
