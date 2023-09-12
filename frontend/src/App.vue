<script>

export default {
  name: 'App',
}

</script>

<template>

  <div id="title">
    <h1>StreamHard</h1>
  </div>
  <div id="content">
    <router-view></router-view>
  </div>


</template>

<style>

@import url('https://fonts.googleapis.com/css2?family=Lato&family=Poppins&family=Satisfy&display=swap');

* {
  box-sizing: border-box;
  padding: 0px;
  margin: 0px;
}

html {
  height: 100vh;
  background-image: url("./assets/photo.png");
  background-size: cover;
  background-repeat: no-repeat;
  background-position: center center;
}

div#app {
  height: 100vh;
}

div#title {
  display: flex;
  justify-content: center;
  align-items: center;
  background: rgb(91, 95, 103);
  color: white;
  height: 3.8rem;
  font-size: 1.05em;
  font-family: 'Satisfy', cursive;
}

div#content {
  display: flex;
  flex-direction: column;
  justify-content: center;
  align-items: center;
  padding: 2rem;
  height: 90%;
}
</style>