import axios from "axios";
import { defineStore } from "pinia";
import { ref } from "vue";

export const useApplicationListStore = defineStore("applicationList", () => {
  // Using any[] for now as the application object structure is complex and not fully defined here
  const applicationList = ref<any[]>([]);

  const fetchApplicationList = async (params: Record<string, any>) => {
    const cleanParams = { ...params };
    Object.keys(cleanParams).forEach(key => {
      if (
        cleanParams[key] === null ||
        cleanParams[key] === undefined ||
        cleanParams[key] === ""
      ) {
        delete cleanParams[key];
      }
    });
    const response = await axios.get("/api/applications", {
      params: cleanParams
    });
    applicationList.value = response.data;
  };

  return {
    applicationList,
    fetchApplicationList
  };
});
