<!-- 申請書ページ描画画面制御 -->
<!-- コンポーネント配置とaxiosで受け取ったjsonを分割し各コンポーネントに渡す-->
<template>
  <div class="application-detail">
    <h1>This is an application detail pagea</h1>
    <div>ID: {{ $route.params.id }}</div>
    <!-- <div>{{ this.application_detail_paper.title }}</div> -->
    <!-- <application-paper :detail="application_detail_paper" :test="test" />
    <application-logs :log="applitaion_detail_log" /> -->
    <application-paper />
    <application-logs />
  </div>
</template>

<script>
import ApplicationPaper from "./components/ApplicationDetail";
import ApplicationLogs from "./components/ApplicationDetailLogs";
export default {
  name: "AppDetail",
  components: {
    ApplicationPaper,
    ApplicationLogs
  },
  data() {
    return {
      //とりあえず申請書詳細はstoreで管理します。
      // test: "piyo",
      // application_detail_paper: {
      //   application_id: "3fa85f64-5717-4562-b3fc-2c963f66afa6",
      //   created_at: "2019-12-13T02:32:53.089Z",
      //   applicant: {
      //     trap_id: "nagatech",
      //     display_name: "ながてち",
      //     is_admin: true
      //   },
      //   current_state: "submitted",
      //   type: "club",
      //   title: "夏コミの交通費をお願いします。",
      //   remarks: "〇〇駅から〇〇駅への移動",
      //   paid_at: "2019-12-13",
      //   ammount: 1200,
      //   repaid_to_id: ["nagatech"],
      //   images: ["3fa85f64-5717-4562-b3fc-2c963f66afa6"]
      // },
      // application_detail_log: {
      //   comments: [
      //     {
      //       comment_id: 0,
      //       user: {
      //         trap_id: "nagatech",
      //         display_name: "ながてち",
      //         is_admin: true
      //       },
      //       comment: "コメント内容",
      //       created_at: "2019-12-13T02:32:53.089Z",
      //       updated_at: "2019-12-13T02:32:53.089Z"
      //     }
      //   ],
      //   state_logs: [
      //     {
      //       update_user: {
      //         trap_id: "nagatech",
      //         display_name: "ながてち",
      //         is_admin: true
      //       },
      //       to_state: "submitted",
      //       reason: "これは雑すぎますね。",
      //       created_at: "2019-12-13T02:32:53.089Z"
      //     }
      //   ],
      //   application_detail_logs: [
      //     {
      //       update_user: {
      //         trap_id: "nagatech",
      //         display_name: "ながてち",
      //         is_admin: true
      //       },
      //       type: "club",
      //       title: "夏コミの交通費をお願いします。",
      //       remarks: "〇〇駅から〇〇駅への移動",
      //       ammount: 1200,
      //       paid_at: "2019-12-13",
      //       updated_at: "2019-12-13T02:32:53.089Z"
      //     }
      //   ],
      //   repayment_logs: [
      //     {
      //       repaid_by_user: {
      //         trap_id: "nagatech",
      //         display_name: "ながてち",
      //         is_admin: true
      //       },
      //       repaid_to_user: {
      //         trap_id: "nagatech",
      //         display_name: "ながてち",
      //         is_admin: true
      //       },
      //       repaid_at: "2019-12-13T02:32:53.089Z"
      //     }
      //   ]
      // }
    };
  }
};
</script>
