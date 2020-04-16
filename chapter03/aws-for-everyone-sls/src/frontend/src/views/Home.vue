<template>
  <div class="ion-page">
    <ion-header>
      <ion-toolbar>
        <ion-title>Add Person</ion-title>
      </ion-toolbar>
    </ion-header>
    <ion-content class="ion-padding">
      <ion-item>
        <ion-input
          :value="firstName"
          @ionInput="firstName = $event.target.value"
          placeholder="Enter first name"
        >
        </ion-input>
      </ion-item>
      <ion-item>
        <ion-input
          :value="lastName"
          @ionInput="lastName = $event.target.value"
          placeholder="Enter last name"
        >
        </ion-input>
      </ion-item>
      <ion-button @click="onClick()" full>Add</ion-button>
    </ion-content>
  </div>
</template>

<script>
import axios from "axios";
const baseUrl = process.env.VUE_APP_API_BASE_URL;
export default {
  name: "home",
  data: function() {
    return {
      firstName: "",
      lastName: ""
    };
  },
  methods: {
    async onClick() {
      console.log(baseUrl);
      console.log(this.firstName);
      console.log(this.lastName);
      try {
        const response = await axios.post(`${baseUrl}persons`, {
          firstName: this.firstName,
          lastName: this.lastName
        });
        console.log(response.data);
        this.firstName = "";
        this.lastName = "";
      } catch (e) {
        console.log(e);
      }
    }
  }
};
</script>
