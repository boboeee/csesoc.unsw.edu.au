<template>
  <v-container class="nav__grid" fluid>
    <v-col>
      <v-row align="center" justify="center">
        <v-card v-for="item in gridItems" :key="item.title" :href="item.link" width="30%" class="ma-3">
          <v-img
            :src="item.src"
            class="white--text"
            height="180px"
            gradient="to bottom, rgba(0,0,0,.3), rgba(0,0,0,.8)"
          >
            <v-card-title class="justify-center fill-height align-center" v-text="item.title"></v-card-title>
          </v-img>
        </v-card>
      </v-row>
    </v-col>
  </v-container>
</template>


<style scoped>
@import url("https://fonts.googleapis.com/css?family=Quicksand");

.nav__grid {
  font-family: "Quicksand", sans-serif;
}
</style>

<script type="text/javascript">
export default {
  name: 'NavGrid',
  // Must be passed from parent object
  // items have title, image url (src), and link
  props: ['gridItems'],

};
</script>
