<template>
  <div class="p-4">
    <div
      v-if="email.htmlBody"
      ref="htmlBody"
      v-html="bodyDecode(email.htmlBody)"
    ></div>
    <div v-else>{{ email.textBody }}</div>
  </div>
</template>

<script>
import quotedPrintable from 'quoted-printable';

const debounce = (fn, ms) => {
  let timer;
  return function () {
    clearTimeout(timer);
    const args = Array.prototype.slice.call(arguments);
    args.unshift(this);
    timer = setTimeout(fn.bind.apply(fn, args), ms);
  };
};

export default {
  props: {
    email: {
      type: Object,
      required: true,
    },
  },
  data() {
    return {
      parents: [],
      scalables: [],
      height: undefined,
      scaler: undefined,
    };
  },
  methods: {
    bodyDecode(s) {
      try {
        return decodeURIComponent(escape(window.atob(s)));
      } catch {
        try {
          return quotedPrintable.decode(s);
        } catch (e) {
          return s;
        }
      }
    },
    findScalables(el, acc) {
      for (let i = 0; i < el.childNodes.length; i++) {
        const child = el.childNodes[i];
        if (child.clientWidth + 48 > document.body.clientWidth) {
          acc.push(child);
          continue;
        }
        this.findScalables(child, acc);
      }
    },
    scale() {
      for (let i = 0; i < this.scalables.length; i++) {
        const el = this.scalables[i];
        el.style.transform = '';
        el.style.transformOrigin = '';
        el.style.position = '';
      }
      for (let i = 0; i < this.parents.length; i++) {
        const el = this.parents[i];
        el.style.height = '';
      }
      setTimeout(() => {
        this.scalables.splice(0);
        this.findScalables(this.$refs.htmlBody, this.scalables);
        const clientWidth = document.body.clientWidth;
        const scrollWidth = document.body.scrollWidth;
        let scale;
        if (clientWidth >= scrollWidth) {
          scale = 1;
        } else {
          scale = (clientWidth - 48) / scrollWidth;
        }
        const parentsHeights = new Map();
        for (let i = 0; i < this.scalables.length; i++) {
          const el = this.scalables[i];
          const height = el.clientHeight * scale;
          el.style.transform = 'scale(' + scale + ')';
          el.style.transformOrigin = 'top left';
          el.style.position = 'absolute';
          const parentHeight = parentsHeights.get(el.parentElement);
          if (parentHeight) {
            parentsHeights.set(el.parentElement, parentHeight + height);
          } else {
            parentsHeights.set(el.parentElement, height);
          }
        }
        this.parents.splice(0);
        for (let [el, height] of parentsHeights) {
          this.parents.push(el);
          el.style.height = height + 'px';
        }
      }, 0);
    },
  },
  mounted() {
    if (this.email.htmlBody) {
      this.scale();
      this.scaler = new ResizeObserver(debounce(this.scale, 100));
      this.scaler.observe(this.$refs.htmlBody);
      const ss = this.$refs.htmlBody.getElementsByTagName('style');
      for (let i = 0; i < ss.length; i++) {
        ss[i].innerHTML = ss[i].innerHTML.replaceAll(/html|body/g, '#email');
      }
    }
  },
  unmounted() {
    if (this.email.htmlBody) {
      this.scaler.disconnect();
    }
  },
};
</script>
