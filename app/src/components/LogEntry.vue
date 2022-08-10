<template>
  <v-card>
    <v-card-title>
      <span class="headline">ISUCON12本選</span>
    </v-card-title>
    <v-data-table
        :headers="headers"
        :items="entries"
        :loading="loading"
        loading-text="Loading... Please wait"
    >
      <template v-slot:item="row">
        <tr>
          <td>{{row.item.timestamp}}</td>
          <td>{{row.item.score}}</td>
          <td>{{row.item.message}}</td>
          <td>
            <v-btn class="mx-2" fab small color="primary" @click="onButtonClick(row.item.access_log)">
              <v-icon dark>mdi-server</v-icon>
            </v-btn>
          </td>
          <td>
            <v-btn class="mx-2" fab dark small color="secondary" @click="onButtonClick(row.item.slow_log)">
              <v-icon dark>mdi-database</v-icon>
            </v-btn>
          </td>
          <td>{{row.item.status}}</td>
        </tr>
      </template>
    </v-data-table>
  </v-card>
</template>

<script>
import axios from "axios";

export default {
  data () {
    return {
      loading: false,
      headers: [
        {
          text: 'Timestamp',
          align: 'start',
          filterable: true,
          value: 'timestamp',
        },
        { text: 'Score', value: 'score' },
        { text: 'Message', value: 'message' },
        { text: 'Access Log', value: 'access_log' },
        { text: 'Slow Log', value: 'slow_log' },
        { text: 'Status', value: 'status' },
      ],
      entries: [],
    }
  },
  methods: {
    getData() {
      this.loading = true;
      return axios
          .get("http://localhost:8082/get?isucon_id=1", {
            dataType: "json",
          })
          .then((response) => {
            console.log(response.data);
            for (let i = 0; i < response.data.length; i++) {
              response.data[i].timestamp = convertTimestamp(response.data[i].timestamp);
              this.entries.push(response.data[i]);
            }
            this.loading = false;
          })
          .catch((err) => alert(err));
    },
    onButtonClick(item) {
      console.log(item);
    },
  },

  mounted() {
    this.getData();
  },
}

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