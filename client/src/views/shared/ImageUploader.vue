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

<script>
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
      images: null,
      uploadImageUrl: [],
      uploadImageBlob: []
    };
  },
  methods: {
    imageChange(files) {
      this.uploadImageUrl = [];
      this.uploadImageBlob = [];
      if (!files) return;
      files.forEach(file => {
        const fr = new FileReader();
        fr.readAsDataURL(file);
        this.uploadImageBlob.push(file);
        this.$emit("update:modelValue", this.uploadImageBlob);
        fr.addEventListener("load", () => {
          this.uploadImageUrl.push(fr.result);
        });
      });
    }
  }
};
</script>
