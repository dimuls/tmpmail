<template>
  <div class="z-10 flex flex-col font-sans shadow">
    <div class="flex flex-col items-center gap-8 p-8 text-white">
      <div class="flex items-center justify-center gap-1 text-xl">
        <mail-icon class="-mb-1 h-7 w-7 -rotate-6" /> Временная почта
      </div>
      <div class="flex gap-2">
        <div
          class="flex cursor-pointer flex-col overflow-hidden whitespace-nowrap rounded-full bg-neutral-800 p-4 shadow active:scale-95"
          @click="copyEmail"
        >
          <div>{{ userEmail }}</div>
          <div class="block text-center text-xs text-neutral-500 xxs:hidden">
            Нажмите чтобы скопировать
          </div>
        </div>
        <button
          class="hidden cursor-pointer rounded-full bg-neutral-800 p-4 shadow active:scale-95 xxs:block"
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
      class="flex w-full flex-wrap items-center justify-center gap-4 bg-white p-4 text-xs"
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
  <div
    class="flex w-full flex-col items-center gap-8 bg-white px-4 py-8 font-sans"
  >
    <div
      v-if="!email"
      class="hidden w-full max-w-3xl overflow-hidden rounded-xl border border-neutral-200 xs:block"
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
          <tr
            v-for="(e, i) in emails"
            :key="'email-' + i"
            class="cursor-pointer"
            @click="email = e"
          >
            <td class="p-4">
              <div>{{ e.from.name }}</div>
              <div class="text-neutral-500">{{ e.from.address }}</div>
            </td>
            <td class="p-4">{{ e.subject }}</td>
            <td class="p-4 text-center">
              <chevron-right-icon class="inline h-5 w-5" />
            </td>
          </tr>
          <tr v-if="loading">
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
    <div v-if="!email" class="flex w-full flex-col gap-4 xs:hidden">
      <div
        v-for="(e, i) in emails"
        :key="'email-' + i"
        class="flex cursor-pointer overflow-hidden rounded-xl border"
        @click="email = e"
      >
        <div class="grid-email grid grow gap-2 p-2">
          <div class="text-right font-bold uppercase text-neutral-300">От</div>
          <div class="grow">
            {{ e.from.name }}
            <span class="text-neutral-500">{{ e.from.address }}</span>
          </div>
          <div class="text-right font-bold uppercase text-neutral-300">
            Тема
          </div>
          <div class="">{{ e.subject }}</div>
        </div>
        <div class="hidden items-center pr-2 xxs:flex">
          <chevron-right-icon class="inline h-5 w-5" />
        </div>
      </div>
      <div class="mt-4 flex items-center justify-center gap-1 text-neutral-500">
        <refresh-icon class="h-4 w-4 animate-spin" /> Ожидание новых писем...
      </div>
    </div>
    <div
      v-if="email"
      ref="emailContainer"
      class="flex w-full max-w-3xl flex-col divide-y rounded-xl border border-neutral-200"
    >
      <div
        class="flex rounded-t-xl bg-neutral-900 p-4 font-sans font-bold uppercase text-white"
      >
        <div
          class="flex cursor-pointer items-center gap-1"
          @click="email = null"
        >
          <chevron-left-icon class="h-5 w-5" /> Назад
        </div>
      </div>
      <div class="flex flex-col gap-4 divide-y p-4 font-sans">
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
      <email :email="email" />
    </div>

    <div v-if="!email" class="flex flex-col gap-8">
      <h1 class="text-center font-sans text-3xl">
        Что такое временная одноразовая почта?
      </h1>
      <p class="w-full max-w-3xl text-center font-sans">
        Временная почта - принимает электронные письма на временный одноразовый
        email, который удаляется через какое-то время. Сервис также известен как
        “почта на 10 минут” или “анонимная почта”. Многие форумы, владельцы
        wi-fi точек, сайты и блоги, требуют от пользователей зарегистрироваться
        до того как они смогут полноценно использовать сайт. Tmp-Mail.ru - это
        новый продвинутый сервис временной почты, который позволит вам навсегда
        забыть о спаме и нежелательной почте.
      </p>
    </div>
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
import Email from './components/Email.vue';
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
    Email,
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
      prolonger: null,
      updater: null,
      loading: true,
      emailContainerWidth: 0,
      emailContainerWidthCalculator: null,
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
  },
  methods: {
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
        this.loading = false;
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

.grid-email {
  grid-template-columns: 3rem auto;
}
</style>

<style scoped>
/*
1. Prevent padding and border from affecting element width. (https://github.com/mozdevs/cssremedy/issues/4)
2. Allow adding a border to an element by just adding a border-width. (https://github.com/tailwindcss/tailwindcss/pull/116)
*/

*,
::before,
::after {
  box-sizing: border-box; /* 1 */
  border-width: 0; /* 2 */
  border-style: solid; /* 2 */
  border-color: theme('borderColor.DEFAULT', currentColor); /* 2 */
}

::before,
::after {
  --tw-content: '';
}

/*
1. Add the correct height in Firefox.
2. Correct the inheritance of border color in Firefox. (https://bugzilla.mozilla.org/show_bug.cgi?id=190655)
3. Ensure horizontal rules are visible by default.
*/

hr {
  height: 0; /* 1 */
  color: inherit; /* 2 */
  border-top-width: 1px; /* 3 */
}

/*
Add the correct text decoration in Chrome, Edge, and Safari.
*/

abbr:where([title]) {
  text-decoration: underline dotted;
}

/*
Remove the default font size and weight for headings.
*/

h1,
h2,
h3,
h4,
h5,
h6 {
  font-size: inherit;
  font-weight: inherit;
}

/*
Reset links to optimize for opt-in styling instead of opt-out.
*/

a {
  color: inherit;
  text-decoration: inherit;
}

/*
Add the correct font weight in Edge and Safari.
*/

b,
strong {
  font-weight: bolder;
}

/*
1. Use the user's configured `mono` font family by default.
2. Correct the odd `em` font sizing in all browsers.
*/

/*
Add the correct font size in all browsers.
*/

small {
  font-size: 80%;
}

/*
Prevent `sub` and `sup` elements from affecting the line height in all browsers.
*/

sub,
sup {
  font-size: 75%;
  line-height: 0;
  position: relative;
  vertical-align: baseline;
}

sub {
  bottom: -0.25em;
}

sup {
  top: -0.5em;
}

/*
1. Remove text indentation from table contents in Chrome and Safari. (https://bugs.chromium.org/p/chromium/issues/detail?id=999088, https://bugs.webkit.org/show_bug.cgi?id=201297)
2. Correct table border color inheritance in all Chrome and Safari. (https://bugs.chromium.org/p/chromium/issues/detail?id=935729, https://bugs.webkit.org/show_bug.cgi?id=195016)
3. Remove gaps between table borders by default.
*/

table {
  text-indent: 0; /* 1 */
  border-color: inherit; /* 2 */
  border-collapse: collapse; /* 3 */
}

/*
1. Change the font styles in all browsers.
2. Remove the margin in Firefox and Safari.
3. Remove default padding in all browsers.
*/

button,
input,
optgroup,
select,
textarea {
  font-family: inherit; /* 1 */
  font-size: 100%; /* 1 */
  font-weight: inherit; /* 1 */
  line-height: inherit; /* 1 */
  color: inherit; /* 1 */
  margin: 0; /* 2 */
  padding: 0; /* 3 */
}

/*
Remove the inheritance of text transform in Edge and Firefox.
*/

button,
select {
  text-transform: none;
}

/*
1. Correct the inability to style clickable types in iOS and Safari.
2. Remove default button styles.
*/

button,
[type='button'],
[type='reset'],
[type='submit'] {
  -webkit-appearance: button; /* 1 */
  background-color: transparent; /* 2 */
  background-image: none; /* 2 */
}

/*
Use the modern Firefox focus style for all focusable elements.
*/

:-moz-focusring {
  outline: auto;
}

/*
Remove the additional `:invalid` styles in Firefox. (https://github.com/mozilla/gecko-dev/blob/2f9eacd9d3d995c937b4251a5557d95d494c9be1/layout/style/res/forms.css#L728-L737)
*/

:-moz-ui-invalid {
  box-shadow: none;
}

/*
Add the correct vertical alignment in Chrome and Firefox.
*/

progress {
  vertical-align: baseline;
}

/*
Correct the cursor style of increment and decrement buttons in Safari.
*/

::-webkit-inner-spin-button,
::-webkit-outer-spin-button {
  height: auto;
}

/*
1. Correct the odd appearance in Chrome and Safari.
2. Correct the outline style in Safari.
*/

[type='search'] {
  -webkit-appearance: textfield; /* 1 */
  outline-offset: -2px; /* 2 */
}

/*
Remove the inner padding in Chrome and Safari on macOS.
*/

::-webkit-search-decoration {
  -webkit-appearance: none;
}

/*
1. Correct the inability to style clickable types in iOS and Safari.
2. Change font properties to `inherit` in Safari.
*/

::-webkit-file-upload-button {
  -webkit-appearance: button; /* 1 */
  font: inherit; /* 2 */
}

/*
Add the correct display in Chrome and Safari.
*/

summary {
  display: list-item;
}

/*
Removes the default spacing and border for appropriate elements.
*/

blockquote,
dl,
dd,
h1,
h2,
h3,
h4,
h5,
h6,
hr,
figure,
p,
pre {
  margin: 0;
}

fieldset {
  margin: 0;
  padding: 0;
}

legend {
  padding: 0;
}

ol,
ul,
menu {
  list-style: none;
  margin: 0;
  padding: 0;
}

/*
Prevent resizing textareas horizontally by default.
*/

textarea {
  resize: vertical;
}

/*
1. Reset the default placeholder opacity in Firefox. (https://github.com/tailwindlabs/tailwindcss/issues/3300)
2. Set the default placeholder color to the user's configured gray 400 color.
*/

input::placeholder,
textarea::placeholder {
  opacity: 1; /* 1 */
  color: theme('colors.gray.400', #9ca3af); /* 2 */
}

/*
Set the default cursor for buttons.
*/

button,
[role='button'] {
  cursor: pointer;
}

/*
Make sure disabled buttons don't get the pointer cursor.
*/
:disabled {
  cursor: default;
}

/*
1. Make replaced elements `display: block` by default. (https://github.com/mozdevs/cssremedy/issues/14)
2. Add `vertical-align: middle` to align replaced elements more sensibly by default. (https://github.com/jensimmons/cssremedy/issues/14#issuecomment-634934210)
   This can trigger a poorly considered lint error in some tools but is included by design.
*/

img,
svg,
video,
canvas,
audio,
iframe,
embed,
object {
  display: block; /* 1 */
  vertical-align: middle; /* 2 */
}

/*
Constrain images and videos to the parent width and preserve their intrinsic aspect ratio. (https://github.com/mozdevs/cssremedy/issues/14)
*/

img,
video {
  max-width: 100%;
  height: auto;
}
</style>
