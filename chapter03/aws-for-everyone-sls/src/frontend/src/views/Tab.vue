<template>
  <div>
    <ion-tabs>
      <ion-tab tab="home">
        <Home />
      </ion-tab>
      <ion-tab tab="persons">
        <Persons />
      </ion-tab>
      <ion-tab-bar slot="bottom">
        <ion-tab-button tab="home">
          <ion-icon name="person-add"></ion-icon>
          <ion-label>Add Person</ion-label>
        </ion-tab-button>
        <ion-tab-button tab="persons">
          <ion-icon name="people"></ion-icon>
          <ion-label>Persons</ion-label>
        </ion-tab-button>
      </ion-tab-bar>
    </ion-tabs>
  </div>
</template>

<script>
const Home = () => import("@/views/Home.vue");
const Persons = () => import("@/views/Persons.vue");
export default {
  name: "tab",
  components: { Home, Persons }
};
</script>
