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

<script lang="ts">
export default {
  props: {
    modelValue: {
      type: Array,
      default: () => []
    }
  },
  emits: ["update:modelValue"],
  data() {
    return {
      images: [] as File[],
      uploadImageUrl: [] as string[],
      uploadImageBlob: [] as File[]
    };
  },
  methods: {
    imageChange(files: File | File[]) {
      this.uploadImageUrl = [];
      this.uploadImageBlob = [];
      if (!files) return;
      const fileList = Array.isArray(files) ? files : [files];
      fileList.forEach((file: File) => {
        const fr = new FileReader();
        fr.readAsDataURL(file);
        this.uploadImageBlob.push(file);
        this.$emit("update:modelValue", this.uploadImageBlob);
        fr.addEventListener("load", () => {
          this.uploadImageUrl.push(fr.result as string);
        });
      });
    }
  }
};
</script>
