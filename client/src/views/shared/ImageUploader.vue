<template>
  <span>
    <v-file-input
      v-model="images"
      label="画像"
      filled
      multiple
      clearable
      chips
      accept="image/*"
      placeholder="画像を添付"
      @change="imageChange"
    ></v-file-input>
    <div v-for="(imageUrl, index) in uploadImageUrl" :key="index">
      <v-img :src="imageUrl" max-width="50%" />
    </div>
  </span>
</template>

<script>
export default {
  data() {
    return {
      images: null,
      uploadImageUrl: [],
      uploadImageBlob: []
    };
  },
  props: {
    value: {
      type: Array,
      value: []
    }
  },
  methods: {
    imageChange(files) {
      this.uploadImageUrl = [];
      this.uploadImageBlob = [];
      files.forEach(file => {
        const fr = new FileReader();
        fr.readAsDataURL(file);
        this.uploadImageBlob.push(file);
        this.$emit("input", this.uploadImageBlob);
        fr.addEventListener("load", () => {
          this.uploadImageUrl.push(fr.result);
        });
      });
    }
  }
};
</script>
