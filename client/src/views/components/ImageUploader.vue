<template>
  <span>
    <v-file-input
      v-model="images"
      multiple
      clearable
      chips
      accept="image/*"
      placeholder="画像を添付"
      @change="imageChange"
    ></v-file-input>
    <img
      v-for="(imageUrl, index) in uploadImageUrl"
      :key="index"
      :src="imageUrl"
    />
  </span>
</template>

<script>
export default {
  data() {
    return {
      images: null,
      uploadImageUrl: []
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
      this.uploadImageBin = [];
      files.forEach(file => {
        const fr = new FileReader();
        const fr2 = new FileReader();
        fr2.readAsDataURL(file);
        fr2.addEventListener("load", () => {
          this.uploadImageUrl.push(fr2.result);
        });
        fr.readAsBinaryString(file);
        fr.addEventListener("load", () => {
          this.uploadImageBin.push(fr.result);
          this.$emit("input", this.uploadImageBin);
        });
      });
    }
  }
};
</script>
