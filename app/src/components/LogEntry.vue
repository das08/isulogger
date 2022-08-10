<template>
  <table-lite
      :has-checkbox="false"
      :is-loading="table.isLoading"
      :is-re-search="table.isReSearch"
      :columns="table.columns"
      :rows="table.rows"
      :rowClasses="table.rowClasses"
      :total="table.totalRecordCount"
      :sortable="table.sortable"
      :messages="table.messages"
      @do-search="doSearch"
      @is-finished="tableLoadingFinish"
      @return-checked-rows="updateCheckedRows"
  ></table-lite>
</template>

<script>
import { defineComponent, reactive } from "vue";
import { useRouter } from 'vue-router'
import TableLite from "vue3-table-lite";
import axios from "axios";

// Fake Data for 'asc' sortable
// const sampleData1 = (offst, limit) => {
//   offst = offst + 1;
//   let data = [];
//   for (let i = offst; i <= limit; i++) {
//     data.push({
//       timestamp: new Date( 2022, 8, 2, i%60, 15, 30 ),
//       score: i*1000+60,
//       message: "test" + i,
//     });
//   }
//   return data;
// };

// Fake Data for 'desc' sortable
// const sampleData2 = (offst, limit) => {
//   let data = [];
//   for (let i = limit; i > offst; i--) {
//     data.push({
//       timestamp: new Date( 2022, 8, 2, i%60, 15, 30 ),
//       score: i*1000+60,
//       message: "test" + i,
//     });
//   }
//   return data;
// };

export default defineComponent({
  name: "App",
  components: { TableLite },
  setup() {
    // Table config
    const table = reactive({
      isLoading: false,
      isReSearch: false,
      // rowClasses: (row) => {
      //   if (row.id == 1) {
      //     return ["aaa", "is-id-one"];
      //   }
      //   return ["bbb", "other"];
      // },
      columns: [
        {
          label: "Timestamp",
          field: "timestamp",
          width: "3%",
          sortable: true,
          isKey: true,
          display: function (row) {
            return (
                convertTimestamp(row.timestamp)
            );
          },
        },
        {
          label: "Score",
          field: "score",
          width: "10%",
          sortable: true,
        },
        {
          label: "Message",
          field: "message",
          width: "15%",
          sortable: true,
        },
        {
          label: "Access Log",
          field: "access_log",
          width: "10%",
          display: function (row) {
            return (
                '<button type="button" data-id="' +
                row.id +
                '" class="is-rows-el quick-btn">Button</button>'
            );
          },
        },
        {
          label: "Slow Log",
          field: "slow_log",
          width: "10%",
          display: function () {
            return (
                '<button @click="moveNextScreen">next</button>'
            );
          },
        },
        {
          label: "Status",
          field: "status",
          width: "5%",
          sortable: false,
        },
      ],
      rows: [],
      totalRecordCount: 0,
      sortable: {
        order: "timestamp",
        sort: "desc",
      },
      messages: {
        pagingInfo: "Showing {0}-{1} of {2}",
        pageSizeChangeLabel: "Row count:",
        gotoPageLabel: "Go to page:",
        noDataAvailable: "No data",
      },
    });

    /**
     * Search Event
     */
    const doSearch = (offset, limit, order, sort) => {
      table.isLoading = true;
      let url = "http://localhost:8082/get?isucon_id=1&sort=" + sort;
      axios.get(url)
          .then((response) => {
            console.log(response.data);

            // refresh table rows
            table.rows = response.data;
            table.totalRecordCount = response.data.length;
            table.sortable.order = order;
            table.sortable.sort = sort;
          });
    };

    /**
     * Loading finish event
     */
    const tableLoadingFinish = (elements) => {
      table.isLoading = false;
      Array.prototype.forEach.call(elements, function (element) {
        if (element.classList.contains("name-btn")) {
          element.addEventListener("click", function () {
            console.log(this.dataset.id + " name-btn click!!");
          });
        }
        if (element.classList.contains("quick-btn")) {
          element.addEventListener("click", function () {
            console.log(this.dataset.id + " quick-btn click!!");
          });
        }
      });
    };

    /**
     * Row checked event
     */
    const updateCheckedRows = (rowsKey) => {
      console.log(rowsKey);
    };

    // First get data
    doSearch(0, 10, "timestamp", "desc");

    const router = useRouter()

    const moveNextScreen = () => {
      router.push('/second')
    }


    return {
      table,
      doSearch,
      tableLoadingFinish,
      updateCheckedRows,
      moveNextScreen,
    };
  },
});

// convert timestamp to date
function convertTimestamp(timestamp) {
  let parsed = Date.parse(timestamp);
  let date = new Date(parsed);
  var year = date.getFullYear();
  var month = date.getMonth();
  var day = date.getDate();
  var hour = date.getHours();
  var min = date.getMinutes();
  var sec = date.getSeconds();
  return year + "/" + month + "/" + day + " " + hour + ":" + min + ":" + sec;
}
</script>

<style scoped>
::v-deep(.vtl-table .vtl-thead .vtl-thead-th) {
  color: #fff;
  background-color: #42b983;
  border-color: #42b983;
}
::v-deep(.vtl-table td),
::v-deep(.vtl-table tr) {
  border: none;
}
::v-deep(.vtl-paging-info) {
  color: rgb(172, 0, 0);
}
::v-deep(.vtl-paging-count-label),
::v-deep(.vtl-paging-page-label) {
  color: rgb(172, 0, 0);
}

</style>