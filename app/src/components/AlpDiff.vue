<template>
  <div class="ranking-wrapper">
    <template v-if="ranking1">
      <v-data-table
        :headers="headers"
        :items="ranking1"
        item-key="uri"
        class="ranking-table ranking-1"
        disable-pagination
        disable-sort
        hide-default-footer
      >

      </v-data-table>
    </template>
    <div>
      <canvas id="rankingLinkCanvas"></canvas>
      <!-- <v-icon>mdi-arrow-right-bold</v-icon> -->
    </div>
    <template v-if="ranking2">
      <v-data-table
        :headers="headers"
        :items="ranking2"
        item-key="uri"
        class="ranking-table ranking-2"
        disable-pagination
        disable-sort
        hide-default-footer
      >
      </v-data-table>
    </template>
  </div>
</template>

<script>
import axios from 'axios'

function makeRanking(alp) {
  if (alp === null) {
    return null;
  }
  const copy = [...alp];
  return copy.sort((a, b) => {
    return b.avg - a.avg;
  }).map((item) => {
    return {
      count: item.count,
      count100: item.count100,
      count200: item.count200,
      count300: item.count300,
      count400: item.count400,
      count500: item.count500,
      uri: item.uri,
      avg: item.avg,
      sum: item.sum
    }
  });
}

export default {
  name: "AlpDiff",
  props: {
    "cmp1": Number,
    "cmp2": Number,
  },
  data() {
    return {
      "alp1": null,
      "alp2": null,
      "headers": [
        {text: "URI", value:"uri"},
        {text:"Avg",value:"avg"},
        {text:"Sum",value:"sum"},
        {text:"Count",value:"count"},
      ],
    }
  },
  computed: {
    ranking1() {
      return makeRanking(this.alp1);
    },
    ranking2() {
      return makeRanking(this.alp2);
    }
  },
  mounted() {
    const pro1 = axios.get("/parsed/alp/" + this.cmp1)
    .then((response) => {
      this.alp1 = response.data;
    })
    const pro2 = axios.get("/parsed/alp/" + this.cmp2)
    .then((response) => {
      this.alp2 = response.data;
    })
    Promise.all([pro1, pro2]).then(() => {
      this.drawRankingLink();
    });

    this.initCanvas();
  },
  methods: {
    initCanvas() {
      const canvas = document.getElementById('rankingLinkCanvas');
      // const ctx = canvas.getContext('2d');
      canvas.style.width = '100%';
      canvas.style.height = '100%';
      canvas.width = canvas.offsetWidth;
      canvas.height = canvas.offsetHeight;
    },
    drawRankingLink() {
      const table1 = document.getElementsByClassName('ranking-1')[0];
      const table2 = document.getElementsByClassName('ranking-2')[0];
      const table1Map = new Map();
      const table2Map = new Map();
      const table1Set = new Set();
      const table2Set = new Set();
      table1.querySelectorAll('td').forEach((td) => {
        const uri = td.innerText;
        if (uri.startsWith("/")) {
          const b = td.getBoundingClientRect();
          const height = (b.top + b.bottom)/2;
          table1Map.set(uri, height);
          table1Set.add(uri);
        }
      });
      table2.querySelectorAll('td').forEach((td) => {
        const uri = td.innerText;
        if (uri.startsWith("/")) {
          const b = td.getBoundingClientRect();
          const height = (b.top + b.bottom)/2;
          table2Map.set(uri, height);
          table2Set.add(uri);
        }
      });

      const uriIntersection = new Set([...table1Set].filter((x) => table2Set.has(x)));

      const canvas = document.getElementById('rankingLinkCanvas');
      const ctx = canvas.getContext('2d');
      ctx.clearRect(0, 0, canvas.width, canvas.height);
      // ctx.fillStyle = 'red';
      // ctx.fillRect(0, 0, canvas.width, canvas.height);
ctx.fillStyle='blue';
        ctx.fillRect(0, 0, 100, 10);
      for (const uri of uriIntersection.keys()) {
        const h1 = table1Map.get(uri);
        const h2 = table2Map.get(uri);
        console.log(uri, h1, h2);
        
      }
    }
  }
}
</script>

<style scoped>
.ranking-wrapper {
  display: grid;
  grid-template-columns: 4fr 1fr 4fr;
}
</style>