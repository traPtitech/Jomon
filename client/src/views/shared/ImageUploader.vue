<template>
  <div>
    <v-file-input
      v-model="images"
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
import { ref } from "vue";

const images = defineModel<File[]>({ default: () => [] });
const uploadImageUrl = ref<string[]>([]);
const uploadImageBlob = ref<File[]>([]);

const imageChange = (files: File | File[]) => {
  uploadImageUrl.value = [];
  uploadImageBlob.value = [];
  if (!files) return;
  const fileList = Array.isArray(files) ? files : [files];
  fileList.forEach((file: File) => {
    const fr = new FileReader();
    fr.readAsDataURL(file);
    uploadImageBlob.value.push(file);
    images.value = uploadImageBlob.value;
    fr.addEventListener("load", () => {
      uploadImageUrl.value.push(fr.result as string);
    });
  });
};
</script>
