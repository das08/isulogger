<template>
  <div class="ranking-wrapper">
    <template v-if="ranking2">
      <v-data-table
        :headers="headers"
        :items="rankingComparison"
        item-key="uri"
        class="ranking-table"
        disable-pagination
        disable-sort
        hide-default-footer
      >
      </v-data-table>
    </template>
  </div>
</template>

<script>
import axios from "axios";

function makeRanking(alp) {
  if (alp === null) {
    return null;
  }
  const copy = [...alp];
  return copy
    .sort((a, b) => {
      return b.avg - a.avg;
    })
    .map((item, i) => {
      return {
        ranking: i + 1,
        count: item.count,
        count100: item.count100,
        count200: item.count200,
        count300: item.count300,
        count400: item.count400,
        count500: item.count500,
        uri: item.uri,
        avg: item.avg,
        sum: item.sum,
      };
    });
}

function makeRankingComparison(alp, baseAlp) {
  const baseUriMap = new Map();
  for (const item of baseAlp) {
    baseUriMap.set(item.uri, item);
  }

  return alp.map((item) => {
    const uri = item.uri;
    const oldItem = baseUriMap.get(uri);
    if (oldItem === undefined) {
      return item;
    }

    return {
      ranking: `${oldItem.ranking} -> ${item.ranking}`,
      count: `${oldItem.count} -> ${item.count}`,
      count100: `${oldItem.count100} -> ${item.count100}`,
      count200: `${oldItem.count200} -> ${item.count200}`,
      count300: `${oldItem.count300} -> ${item.count300}`,
      count400: `${oldItem.count400} -> ${item.count400}`,
      count500: `${oldItem.count500} -> ${item.count500}`,
      uri: item.uri,
      avg: `${oldItem.avg} -> ${item.avg}`,
      sum: `${oldItem.sum} -> ${item.sum}`,
    };
  });
}

export default {
  name: "AlpDiff",
  props: {
    cmp1: Number,
    cmp2: Number,
  },
  data() {
    return {
      alp1: null,
      alp2: null,
      headers: [
        { text: "#", value: "ranking" },
        { text: "URI", value: "uri" },
        { text: "Avg", value: "avg" },
        { text: "Sum", value: "sum" },
        { text: "Count", value: "count" },
      ],
    };
  },
  computed: {
    ranking1() {
      return makeRanking(this.alp1);
    },
    ranking2() {
      return makeRanking(this.alp2);
    },
    rankingComparison() {
      if (this.ranking1 === null || this.ranking2 === null) {
        return null;
      }
      return makeRankingComparison(this.ranking2, this.ranking1);
    },
  },
  mounted() {
    axios.get("/parsed/alp/" + this.cmp1).then((response) => {
      this.alp1 = response.data;
    });
    axios.get("/parsed/alp/" + this.cmp2).then((response) => {
      this.alp2 = response.data;
    });
  },
  methods: {},
};
</script>

<style scoped>
/* .ranking-wrapper {
  display: grid;
  grid-template-columns: 4fr 1fr 4fr;
} */
</style>
