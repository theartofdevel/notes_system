<template>
  <form class="form" @submit.prevent="submitForm">
    <h1 class="heading">NS</h1>
    <input type="email" v-model="email" autocomplete="email" class="input" placeholder="Введите електронную почту"/>
    <div class="password__container">
      <input autocomplete="current-password" :type="passwordInputType"
             class="input" v-model="password" placeholder="Введите пароль"/>
      <button type="button" @click="togglePasswordVisible" class="password__button">
        <eye v-if="passwordInputType === 'password'" style="font-size: 25px"/>
        <eye-off v-else style="font-size: 25px"/>
      </button>
    </div>
    <button type="submit" class="button">Войти</button>
  </form>
</template>

<script lang="ts">
import {defineComponent} from "vue";
import {Eye, EyeOff} from 'mdue'

export default defineComponent({
  components: {Eye, EyeOff},
  data: () => ({
    passwordInputType: "password",
    email: '',
    password: '',
  }),
  methods: {
    togglePasswordVisible() {
      if (this.$data.passwordInputType === 'password') {
        this.$data.passwordInputType = 'text';
      } else {
        this.$data.passwordInputType = 'password';
      }
    },
    submitForm() {
      console.info(`SUBMITTING FORM WITH DATA : ${JSON.stringify({
        email: this.$data.email,
        password: this.$data.password
      })}`)
    }
  }
})
</script>

<style>
body {
  height: 100vh;
  width: 100%;
  overflow: hidden;
  display: flex;
  justify-content: center;
  align-items: center;
}

.form {
  padding: 100px;
  border: 3px solid rgba(153, 204, 255, 1);
  border-radius: 10px;
  display: flex;
  flex-direction: column;
}

.heading {
  text-align: center;
  font-size: 70px;
  color: dodgerblue;
  font-family: 'Pacifico', cursive;
}

.input {
  outline: none;
  border: 3px solid lightblue;
  border-radius: 10px;
  padding: 10px 15px;
  font-weight: 500;
  margin: 10px 0;
}

.input:focus-visible {
  outline-offset: 0;
  outline: darkblue auto 1px;
}

.button {
  outline: none;
  padding: 10px;
  margin: 10px 0;
  background-color: dodgerblue;
  border: 3px solid dodgerblue;
  color: #fff;
  font-weight: 500;
  border-radius: 10px;
}

.button:hover {
  cursor: pointer;
}

.button:focus-visible {
  background-color: darkblue;
  border-color: darkblue;
}

.password__container {
  display: flex;
  align-items: center;
  position: relative
}

.password__button {
  outline: none;
  border: none;
  display: flex;
  background: none;
  justify-content: center;
  align-items: center;
  position: absolute;
  top: 0;
  bottom: 0;
  right: 15px;
}

.password__button:focus-visible {
  color: darkblue;
}

.password__button:hover {
  cursor: pointer;
}
</style>