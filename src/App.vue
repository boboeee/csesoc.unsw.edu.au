<template>
  <v-app id="main-app">
    <!-- Navigation bar/app bar goes here -->
    <v-content>
<div>
      <Sidebar :drawer="drawer" />
      <v-app-bar app dark>
        <v-app-bar-nav-icon class="ma-2" @click.stop="drawer = !drawer"> </v-app-bar-nav-icon>
        <div class="flex-grow-1"></div>

        <router-link to="/"
          ><v-container
            class="fill-height"
            fluid
            style="max-height: 64px; max-width:100px"
          >
            <v-img class="ma-8" src="https://github.com/csesoc/csesoc.unsw.edu.au/blob/frontendCombined/src/assets/csesoclogobluewhite.png?raw=true" />
        </v-container>
        </router-link>
         <div class="flex-grow-1"></div>
        <LoginForm align="right"/>
      </v-app-bar>

  <router-view></router-view>
    <Footer/>
  </div>
    </v-content>
  </v-app>
</template>

<script>
import Footer from '@/components/Footer';
import Sidebar from '@/components/Sidebar.vue';
import LoginForm from '@/components/LoginForm';

export default {
  name: 'App',
  data: () => ({
    drawer: false,
    gridApiUri: 'https://gistcdn.githack.com/gawdn/464b5ed74404481f7296fb24f9f28243/raw/c9f63e5117a1406db9af5266c8cfd448161bbfec/test_grid.json',
    gridItems: [],
    listApiUri: 'https://gistcdn.githack.com/gawdn/464b5ed74404481f7296fb24f9f28243/raw/c9f63e5117a1406db9af5266c8cfd448161bbfec/test_grid.json',
    listItems: [],
  }),

  components: {
    Footer,
    Sidebar,
    LoginForm
  },

  mounted() {
    fetch(this.gridApiUri)
      .then(r => r.json())
      .then((responseJson) => {
        this.gridItems = responseJson;
      });
  },

};
</script>

<style>
    @import url("https://fonts.googleapis.com/css?family=Quicksand&display=swap");
* {
  margin: 0;
  padding: 0;
}
#showcase {
  background-image: url("https://backgroundcheckall.com/wp-content/uploads/2017/12/black-tech-background-12.jpg");
  background-size: cover;
  background-position: center;
  height: 100vh;
  display: flex;
  flex-direction: column;
  justify-content: center;
  align-items: center;
  text-align: center;
  padding: 0 20px;
}
#showcase img {
  max-width: 30%;
  max-height: 30vh;
}
#showcase h1 {
  font-size: 50px;
  line-height: 1.2;
}
#showcase p {
  font-size: 20px;
}
#showcase .button {
  font-family: "Quicksand", sans-serif;
  font-size: 18px;
  text-decoration: none;
  color: #fff;
  background: #3a76f8;
  padding: 10px 20px;
  border-radius: 10px;
  margin-top: 20px;
  width: 250px;
}
#showcase .button:hover {
  background: #3a76f8;
  color: #fff;
}
.content {
    padding: 2vw 15vw 5vw 15vw;
}
</style>
