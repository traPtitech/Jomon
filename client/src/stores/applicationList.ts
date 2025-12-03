import { ApplicationList } from "@/types/application";
import axios from "axios";
import { defineStore } from "pinia";
import { ref } from "vue";

export const useApplicationListStore = defineStore("applicationList", () => {
  const applicationList = ref<ApplicationList[]>([]);

  const fetchApplicationList = async (params: Record<string, unknown>) => {
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
    const response = await axios.get<ApplicationList[]>("/api/applications", {
      params: cleanParams
    });
    applicationList.value = response.data;
  };

  return {
    applicationList,
    fetchApplicationList
  };
});
