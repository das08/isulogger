<template>
  <v-card>
    <v-dialog v-model="showModal" max-width="500px">
      <v-card>
        <v-card-title class="headline">Message</v-card-title>
        <v-card-text>
          <div contenteditable="true" v-html="selectedMessage" @input="onMessageInput" aria-multiline="true" role="textbox">
          </div>
        </v-card-text>
        <v-card-actions>
          <v-spacer></v-spacer>
          <v-btn color="primary" text @click="onCloseMessageModal">Close</v-btn>
        </v-card-actions>
      </v-card>
    </v-dialog>

    <v-alert v-model="error_alert" dense text type="error">
      {{ error_message }}
    </v-alert>
    <v-card-title>
      <v-container fluid class="headline">
        <v-row align="center">
          <v-col
              class="d-flex"
              cols="12"
              sm="6"
          >
            <v-select
                v-model="selected_contest"
                @change="onContestSelect"
                :items="contests"
                :item-text="item =>`${item.contest_name} (Contest ID: ${item.contest_id})`"
                item-value="contest_id"
                label="選択中のコンテスト"
                outlined
            ></v-select>
          </v-col>
        </v-row>
      </v-container>

<!--      <span class="headline">ISUCON12本選</span>-->
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

      <template v-slot:[`item.message`]="{ item }">
        <div @click="onMessageClick(item)">
          {{ item.message }}
        </div>
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

      <template v-slot:[`item.delete_log`]="{ item }">
        <v-btn class="mx-2" icon x-small color="error" @click="onDeleteEntry(item)">
          <v-icon dark>mdi-trash-can</v-icon>
        </v-btn>
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
            <span class="text-h5">【{{ log_type }} Result】 {{ log_contents.timestamp }} (Score: {{ log_contents.score }})</span>
          </v-card-title>
          <v-card-text>
            <span class="text-subtitle-1">{{ log_contents.file_name }}</span>
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
import { authHeaders } from "@/ls";
import axios from "axios";

export default {
  data () {
    return {
      loading: false,
      contests: [],
      selected_contest: null,
      headers: [
        {
          text: 'Timestamp',
          align: 'start',
          filterable: true,
          value: 'timestamp',
          width: '15%',
        },
        { text: 'Best', value: 'max_score', width: '10%', align: 'end'},
        { text: 'Score', value: 'score', width: '10%' },
        { text: 'Message', value: 'message', width: '30%' },
        { text: 'Access Log', value: 'access_log', width: '10%' },
        { text: 'Slow Log', value: 'slow_log', width: '10%' },
        { text: 'Status', value: 'status', width: '5%' },
        { text: 'Branch', value: 'branch_name', width: '10%' },
        { text: '', value: 'delete_log', width: '10%' },
      ],
      entries: [],
      dialog: false,
      log_type: "",
      log_contents: [],
      error_alert: false,
      error_message: "",

      selectedEntryId: -1,
      selectedMessage: "",
      editedMessage: "",
      showModal: false,
    }
  },
  methods: {
    getLogEntry(contestID) {
      if (contestID === undefined) {
        return;
      }
      this.loading = true;
      let entries = [];
      return axios
          .get("/api/entry?contest_id="+contestID, {
            dataType: "json",
            headers: authHeaders(),
          })
          .then((response) => {
            if (response.data === null) {
              this.entries = [];
              this.loading = false;
              return;
            }
            for (let i = 0; i < response.data.length; i++) {
              response.data[i].timestamp = convertTimestamp(response.data[i].timestamp);
              entries.push(response.data[i]);
            }
            this.loading = false;
            this.error_alert = false;
            this.entries = entries;
          })
          .catch((err) => {
            this.error_alert = true;
            this.loading = false;
            this.error_message = err.message;
          });
    },

    getContest() {
      return axios
          .get("/api/contest", {
            dataType: "json",
            headers: authHeaders(),
          })
          .then((response) => {
            if (response.data === null) {
              this.contests = [];
              return;
            }
            this.contests = response.data;
            this.error_alert = false;
          })
          .catch((err) => {
            this.error_alert = true;
            this.loading = false;
            this.error_message = err.message;
          });
    },

    onButtonClick(item, logType) {
      let filePath = "";
      if (logType === "Access Log") {
        filePath = item.access_log_path;
      } else if (logType === "Slow Log") {
        filePath = item.slow_log_path;
      }
      if (filePath === undefined || filePath === "") {
        return;
      }
      return axios
          .get("/api/log/"+filePath+"?id=1", {
            dataType: "text",
            headers: authHeaders(),
          })
          .then((response) => {
            this.log_contents = {
              file_name: filePath,
              log: response.data,
              score: item.score,
              timestamp: item.timestamp,
            }

            this.log_type = logType;
            this.dialog = true;
            this.error_alert = false;
          })
          .catch((err) => {
            this.error_alert = true;
            this.loading = false;
            this.error_message = err.message;
          });
    },

    onMessageClick(item) {
      // console.log('onMessage click: ', item);

      const lines = item.message.split("\n");
      // for each line, create one div
      let message = "";
      for (let i = 0; i < lines.length; i++) {
        message += "<div>" + lines[i] + "</div>";
      }
      this.selectedMessage = message;

      this.editedMessage = item.message;
      this.selectedEntryId = item.id;
      this.showModal = true;
    },

    onMessageInput(event) {
      this.editedMessage = event.target.innerText;
    },

    onCloseMessageModal() {
      this.showModal = false;
      this.selectedMessage = "";

      const selectedEntryId = this.selectedEntryId;
      const editedMessage = this.editedMessage;

      this.selectedEntryId = -1;
      this.editedMessage = "";

      // update message with editedMessage
      axios.put("/api/entry/" + this.selected_contest + "/" +  selectedEntryId + "/message", {
        message: editedMessage,
      }, {
        dataType: "json",
        headers: authHeaders(),
      }).catch((err) => {
        this.error_alert = true;
        this.loading = false;
        this.error_message = err.message;
      }).then(() => {
        // trigger refresh immediately
        this.getLogEntry(this.selected_contest);
      });
    },

    onDeleteEntry(item) {
      this.loading = true;
      return axios
          .delete("/api/entry/"+item.id, {
            dataType: "json",
            headers: authHeaders(),
          })
          .then(() => {
            this.loading = false;
            this.error_alert = false;
            this.getLogEntry(this.selected_contest);
          })
          .catch((err) => {
            this.error_alert = true;
            this.loading = false;
            this.error_message = err.message;
          });
    },

    onContestSelect(contestID) {
      this.selected_contest = contestID;
      localStorage.setItem("contest_id", contestID);
      this.getLogEntry(contestID);
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
    }
  },

  mounted() {
    if(Object.prototype.hasOwnProperty.call(localStorage, "contest_id")) {
      this.selected_contest = parseInt(JSON.parse(JSON.stringify(localStorage.getItem("contest_id"))));
      console.log("mounted",this.selected_contest, typeof this.selected_contest);
      this.getLogEntry(this.selected_contest);
    }
    this.getContest();

    setInterval(() => {
      this.getLogEntry(this.selected_contest);
    }, 5000);
  },
}

// convert timestamp to date
function convertTimestamp(timestamp) {
  let parsed = Date.parse(timestamp);
  let date = new Date(parsed);
  var year = date.getFullYear();
  var month = date.getMonth() + 1;
  var day = date.getDate();
  var hour = ('0'+date.getHours()).slice(-2);
  var min = ('0'+date.getMinutes()).slice(-2);
  var sec = ('0'+date.getSeconds()).slice(-2);
  return year + "/" + month + "/" + day + " " + hour + ":" + min + ":" + sec;
}
</script>

<style scoped>
#log_dialog {
  white-space: pre;
  word-wrap: normal;
  font-family: Menlo,Monaco,Consolas,"Courier New",monospace;
  font-size: 12px;
  color: black;
}
</style>
