<template>
  <v-card>
    <v-card-title>
      <span class="headline">ISUCON12本選</span>
    </v-card-title>
    <v-data-table
        :headers="headers"
        :items="entries"
        :loading="loading"
        :footer-props="{
          'items-per-page-options': [10, 20, 30, 40, 50]
        }"
        :items-per-page="20"
        loading-text="Loading... Please wait"
    >

      <template v-slot:[`item.max_score`]="{ item }">
        <template v-if="item.score === maxScore()">
          <v-icon small class="mr-2"> mdi-crown </v-icon>
        </template>
      </template>

      <template v-slot:[`item.score`]="{ item }">
        <strong>{{ item.score }}</strong>
      </template>

      <template v-slot:[`item.compare1`]="{ item }">
        <input type="radio" name="compare1" @value="item.id" @click="chooseCompare(1, item.id)">
      </template>

      <template v-slot:[`item.compare2`]="{ item }">
        <input type="radio" name="compare2" @value="item.id" @click="chooseCompare(2, item.id)">
      </template>

      <template v-slot:[`item.access_log`]="{ item }">
        <template v-if="item.access_log_path === ''">
          <v-btn class="mx-2" fab small color="primary" disabled>
            <v-icon dark>mdi-server</v-icon>
          </v-btn>
        </template>
        <template v-else>
          <v-btn class="mx-2" fab small color="primary" @click="onButtonClick(item, 'Access Log')">
            <v-icon dark>mdi-server</v-icon>
          </v-btn>
        </template>

      </template>

      <template v-slot:[`item.slow_log`]="{ item }">
        <template v-if="item.slow_log_path === ''">
          <v-btn class="mx-2" fab small color="secondary" disabled>
            <v-icon dark>mdi-database</v-icon>
          </v-btn>
        </template>
        <template v-else>
          <v-btn class="mx-2" fab dark small color="secondary" @click="onButtonClick(item, 'Slow Log')">
            <v-icon dark>mdi-database</v-icon>
          </v-btn>
        </template>

      </template>

      <template v-slot:[`item.status`]="{}">
        <v-chip color="green" outlined small>
          Completed
        </v-chip>
      </template>

    </v-data-table>

    <v-row justify="center">
      <v-dialog
          v-model="dialog"
          scrollable
          width="95%"
      >
        <v-card>
          <v-card-title>
            <span class="text-h5">{{ log_type }} Result</span>
          </v-card-title>
          <v-card-text>
            <span class="text-subtitle-1">{{ log_contents.timestamp }} (Score: {{ log_contents.score }})</span>
          </v-card-text>
          <v-card-text id="log_dialog">
            {{ log_contents.log }}
          </v-card-text>
        </v-card>
      </v-dialog>
    </v-row>
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
          width: '15%',
        },
        { text: 'Best', value: 'max_score', width: '10%', align: 'end' },
        { text: 'Score', value: 'score', width: '10%' },
        { text: 'Message', value: 'message', width: '15%' },
        { text: 'Cmp 1', value: 'compare1', width: '5%' },
        { text: 'Cmp 2', value: 'compare2', width: '5%' },
        { text: 'Access Log', value: 'access_log', width: '10%' },
        { text: 'Slow Log', value: 'slow_log', width: '10%' },
        { text: 'Status', value: 'status', width: '20%' },
      ],
      entries: [],
      dialog: false,
      log_type: "",
      log_contents: [],
      compare: {
        compare1: null,
        compare2: null,
      },
    }
  },
  methods: {
    getData() {
      this.loading = true;
      return axios
          .get("http://localhost:8082/get?contest_id=1", {
            dataType: "json",
          })
          .then((response) => {
            if (response.data === null) {
              this.entries = [];
              this.loading = false;
              return;
            }
            console.log(response.data);
            for (let i = 0; i < response.data.length; i++) {
              response.data[i].timestamp = convertTimestamp(response.data[i].timestamp);
              this.entries.push(response.data[i]);
            }
            this.loading = false;
          })
          .catch((err) => alert(err));
    },

    onButtonClick(item, logType) {
      let filePath = "";
      if (logType === "Access Log") {
        filePath = item.access_log_path;
      } else if (logType === "Slow Log") {
        filePath = item.slow_log_path;
      }
      console.log(filePath);
      if (filePath === undefined || filePath === "") {
        return;
      }
      return axios
          .get("http://localhost:8082/log/"+filePath+"?id=1", {
            dataType: "text",
          })
          .then((response) => {
            console.log(response.data);
            this.log_contents = {
              log: response.data,
              score: item.score,
              timestamp: item.timestamp,
            }

            this.log_type = logType;
            this.dialog = true;
          })
          .catch((err) => alert(err))
    },

    compareScore(index) {
      let maxScore = this.entries.reduce((a,b)=>a.score>b.score?a:b);
      if (this.entries[index].score === maxScore.score) {
        return "blue";
      }
      if (index < this.entries.length -1 ) {
        if (this.entries[index].score > this.entries[index + 1].score) {
          return "green";
        } else if (this.entries[index].score < this.entries[index + 1].score) {
          return "red";
        } else {
          return "black";
        }
      } else {
        return "black";
      }
    },

    maxScore() {
      return this.entries.reduce((a,b)=>a.score>b.score?a:b).score;
    },

    chooseCompare(compareIndex, id) {
      if (compareIndex === 1) {
        this.compare.compare1 = id;
      } else if (compareIndex === 2) {
        this.compare.compare2 = id;
      }
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

<style scoped>
#log_dialog {
  white-space: pre;
  word-wrap: normal;
  font-family: Monaco,monospace;
  font-size: 12px;
}
</style>
