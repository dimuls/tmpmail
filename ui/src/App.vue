<template>
  <div class="z-10 flex flex-col font-sans shadow">
    <div class="flex flex-col items-center gap-8 p-8 text-white">
      <div class="flex items-center justify-center gap-1 text-xl">
        <mail-icon class="-mb-1 h-7 w-7 -rotate-6" /> Временная почта
      </div>
      <div class="flex gap-2">
        <div
          class="overflow-hidden whitespace-nowrap rounded-full bg-neutral-800 p-4 shadow"
        >
          {{ userEmail }}
        </div>
        <button
          class="cursor-pointer rounded-full bg-neutral-800 p-4 shadow active:scale-95"
          @click="copyEmail"
        >
          <clipboard-copy-icon class="h-6 w-6" />
        </button>
      </div>
      <div class="w-full max-w-md text-center text-xs text-neutral-500">
        Временная одноразовая почта - это лучшая защита от спама и нежелательной
        почты. Доступно анонимно, без регистрации и бесплатно.
      </div>
    </div>
    <div
      class="flex w-full items-center justify-center gap-2 bg-white p-4 text-xs"
    >
      <button
        class="flex cursor-pointer gap-1 rounded-full bg-neutral-50 p-4 shadow active:scale-95"
        @click="copyEmail"
      >
        <clipboard-copy-icon class="h-4 w-4" /> Копировать
      </button>
      <button
        class="flex cursor-pointer gap-1 rounded-full bg-neutral-50 p-4 shadow active:scale-95"
        @click="getAccount(true)"
      >
        <refresh-icon class="h-4 w-4" /> Обновить
      </button>
      <button
        class="flex cursor-pointer gap-1 rounded-full bg-neutral-50 p-4 shadow active:scale-95"
        @click="getNewAccount"
      >
        <pencil-alt-icon class="h-4 w-4" /> Сменить
      </button>
      <button
        class="flex cursor-pointer gap-1 rounded-full bg-neutral-50 p-4 shadow active:scale-95"
        @click="getNewAccount"
      >
        <trash-icon class="h-4 w-4" /> Удалить
      </button>
    </div>
  </div>
  <div class="flex w-full flex-col items-center gap-8 bg-white p-8">
    <div
      v-if="!email"
      class="w-full max-w-3xl overflow-hidden rounded-xl border border-neutral-200"
    >
      <table class="w-full table-auto">
        <thead>
          <tr class="bg-neutral-900 font-bold uppercase text-white">
            <td class="p-4">Отправитель</td>
            <td class="grow p-4">Тема</td>
            <td class="p-4">Просмотр</td>
          </tr>
        </thead>
        <tbody>
          <tr v-for="(e, i) in emails" :key="'email-' + i">
            <td class="p-4">
              <div>{{ e.from.name }}</div>
              <div class="text-neutral-500">{{ e.from.address }}</div>
            </td>
            <td class="p-4">{{ e.subject }}</td>
            <td class="cursor-pointer p-4 text-center" @click="email = e">
              <chevron-right-icon class="inline h-5 w-5" />
            </td>
          </tr>
          <tr>
            <td colspan="3" class="p-8">
              <div
                class="flex items-center justify-center gap-1 text-neutral-500"
              >
                <refresh-icon class="h-4 w-4 animate-spin" /> Ожидание новых
                писем...
              </div>
            </td>
          </tr>
        </tbody>
      </table>
    </div>
    <div
      v-if="email"
      class="flex w-full max-w-3xl flex-col divide-y overflow-hidden rounded-xl border border-neutral-200"
    >
      <div class="flex bg-neutral-900 p-4 font-bold uppercase text-white">
        <div
          class="flex cursor-pointer items-center gap-1"
          @click="email = null"
        >
          <chevron-left-icon class="h-5 w-5" /> Назад
        </div>
      </div>
      <div class="flex flex-col gap-4 divide-y p-4">
        <div class="flex justify-between gap-8">
          <div class="flex flex-col">
            <div>{{ email.from.name }}</div>
            <div class="text-neutral-500">{{ email.from.address }}</div>
          </div>
          <div class="flex flex-col">
            <div class="text-right text-neutral-500">Дата и время:</div>
            <div class="text-right">{{ email.date }}</div>
          </div>
        </div>
        <div class="flex gap-4 pt-4">
          <div class="text-neutral-500">Тема:</div>
          <div>{{ email.subject }}</div>
        </div>
      </div>
      <div
        class="relative p-4"
        :style="{
          'min-height': htmlBody.height + 'px',
        }"
      >
        <div
          v-if="email.htmlBody"
          ref="htmlBody"
          class="origin-top-left translate-x-0 translate-y-0"
          :class="{
            absolute: htmlBody.scale != 1,
          }"
          :style="{
            transform: 'scale(' + htmlBody.scale + ')',
          }"
          v-html="base64Decode(email.htmlBody)"
        ></div>
        <div v-if="!email.htmlBody">{{ email.textBody }}</div>
      </div>
    </div>

    <h1 class="text-3xl">Что такое временная одноразовая почта?</h1>
    <p class="w-full max-w-3xl text-center">
      Временная почта - принимает электронные письма на временный одноразовый
      email, который удаляется через какое-то время. Сервис также известен как
      “почта на 10 минут” или “анонимная почта”. Многие форумы, владельцы wi-fi
      точек, сайты и блоги, требуют от пользователей зарегистрироваться до того
      как они смогут полноценно использовать сайт. Tmp-Mail.ru - это новый
      продвинутый сервис временной почты, который позволит вам навсегда забыть о
      спаме и нежелательной почте.
    </p>
  </div>
  <div class="flex flex-col font-sans">
    <div
      class="items-center gap-8 bg-neutral-900 p-8 text-center text-xs text-neutral-500 text-white"
    >
      По вопросам и предложениям пишите на
      <a class="underline" href="https://t.me/dimuls">телеграм</a> или
      <a
        class="underline"
        href="mailto: dimuls@yandex.ru?subject=Вопросы и предложения по tmp-mail.ru"
        >почту</a
      >
    </div>
  </div>
</template>

<script>
import axios from 'axios';
import {
  MailIcon,
  ClipboardCopyIcon,
  RefreshIcon,
  PencilAltIcon,
  ChevronRightIcon,
  ChevronLeftIcon,
  TrashIcon,
} from '@heroicons/vue/solid';
import { mimeWordsDecode } from 'emailjs-mime-codec';
import emailAddrs from 'email-addresses';
import dayjs from 'dayjs';

const domain = 'tmp-mail.ru';
const baseURL = 'https://tmp-mail.ru/api';
const tokenHeader = 'X-TOKEN';
const tokenKey = 'token';

export default {
  components: {
    MailIcon,
    ClipboardCopyIcon,
    RefreshIcon,
    PencilAltIcon,
    ChevronRightIcon,
    ChevronLeftIcon,
    TrashIcon,
  },
  data() {
    return {
      token: null,
      api: axios.create({ baseURL }),
      ttl: null,
      username: null,
      emails: [],
      email: null,
      htmlBody: {
        calculator: null,
        scale: 1,
        height: null,
      },
      prolonger: null,
      updater: null,
    };
  },
  computed: {
    userEmail() {
      if (this.username) {
        return this.username + '@' + domain;
      } else {
        return 'Загрузка...';
      }
    },
  },
  watch: {
    email(e) {
      if (!e) {
        this.htmlBody.calculator.disconnect();
        return;
      }

      setTimeout(() => {
        this.htmlBody.calculator.observe(document.body);
        this.recalculateHTMLBody();
      }, 0);
    },
    async token(t) {
      if (this.prolonger) {
        clearInterval(this.prolonger);
        this.prolonger = null;
      }
      if (this.updater) {
        clearInterval(this.updater);
        this.updater = null;
      }
      if (!t) {
        return;
      }
      this.api = axios.create({
        baseURL,
        headers: {
          [tokenHeader]: t,
        },
      });
      this.prolonger = setInterval(() => this.prolongAccount(), 60000);
      this.updater = setInterval(() => this.getAccount(), 10000);
      await this.getAccount();
    },
  },
  async mounted() {
    this.token = localStorage.getItem(tokenKey);
    if (!this.token) {
      await this.getNewAccount();
    }
    this.htmlBody.calculator = new ResizeObserver(this.recalculateHTMLBody);
  },
  methods: {
    recalculateHTMLBody() {
      if (document.body.offsetWidth - 99 < this.$refs.htmlBody.offsetWidth) {
        this.htmlBody.scale =
          (document.body.offsetWidth - 99) / this.$refs.htmlBody.offsetWidth;
        this.htmlBody.height =
          this.$refs.htmlBody.offsetHeight * this.htmlBody.scale + 32;
      } else {
        this.htmlBody.scale = 1;
      }
    },
    base64Decode(s) {
      try {
        return decodeURIComponent(escape(window.atob(s)));
      } catch {
        s = s.replace(/={1}/g, '%');
        return decodeURIComponent(s);
      }
    },
    mimeWordsDecode(s) {
      return mimeWordsDecode(s);
    },
    async copyEmail() {
      await navigator.clipboard.writeText(this.userEmail);
    },
    async getNewAccount() {
      try {
        await this.api.delete('/account');
      } catch (e) {
        console.error(e);
      }
      try {
        const res = await this.api.post('/account');
        localStorage.setItem(tokenKey, res.data);
        this.email = null;
        this.token = res.data;
      } catch (e) {
        console.error(e);
      }
    },
    async prolongAccount() {
      try {
        await this.api.patch('/account');
      } catch (e) {
        console.error(e);
      }
    },
    async getAccount(resetEmail) {
      try {
        const res = await this.api.get('/account');
        if (resetEmail) {
          this.email = null;
        }
        this.ttl = res.data.ttl;
        this.username = res.data.username;
        this.emails.splice(0);
        if (res.data.emails) {
          this.emails.push(
            ...res.data.emails.map((e) => ({
              ...e,
              date: dayjs(e.date).format('YYYY-MM-DD HH:mm:ss'),
              from: e.from
                .map(this.mimeWordsDecode)
                .map(emailAddrs.parseOneAddress)[0],
            }))
          );
        }
      } catch (e) {
        if (e.response.status === 404) {
          await this.getNewAccount();
        } else {
          console.error(e);
        }
      }
    },
  },
};
</script>

<style>
@tailwind base;
@tailwind components;
@tailwind utilities;

html,
body,
#app {
  margin: 0;
  padding: 0;
  width: 100%;
  min-height: 100vh;
}
#app {
  display: flex;
  flex-direction: column;
  align-items: stretch;
  justify-content: stretch;
  background-color: rgb(23, 23, 23);
}
</style>
