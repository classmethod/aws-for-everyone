<template>
  <div class="ion-page">
    <ion-header>
      <ion-toolbar>
        <ion-title>Persons</ion-title>
      </ion-toolbar>
    </ion-header>
    <ion-content class="ion-padding">
      <ion-list v-bind:key="person.Id" v-for="person in persons">
        <ion-item>
          <ion-label>{{ person.FirstName }} {{ person.LastName }}</ion-label>
          <ion-button @click="deleteUser(person.Id)" full>delete</ion-button>
        </ion-item>
      </ion-list>
    </ion-content>
  </div>
</template>

<script>
import axios from "axios";
const baseUrl = process.env.VUE_APP_API_BASE_URL;
export default {
  name: "persons",
  data() {
    return {
      persons: null
    };
  },
  watch: {
    $route: "reload"
  },
  async created() {
    await this.reload();
  },
  methods: {
    async reload() {
      console.log(baseUrl);
      if (this.$route.fullPath == "/persons") {
        try {
          const response = await axios.get(`${baseUrl}persons`);
          this.persons = response.data;
        } catch (e) {
          console.log(e);
        }
      }
    },
    async deleteUser(id) {
      try {
        await axios.delete(`${baseUrl}persons/${id}`);
      } catch (e) {
        console.log(e);
      }
      await this.reload();
    }
  }
};
</script>
