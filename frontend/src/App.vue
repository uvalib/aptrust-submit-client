<template>
   <Toast position="center" />
   <ConfirmDialog position="top" />
   <header id="aptsubmit-header" class="main-header">
      <div class="library-link">
          <UvaLibraryLogo />
      </div>
      <div class="site-link">
          <router-link to="/">
            APTrust Submissions
          </router-link>
      </div>
   </header>
   <main>
      <RouterView  v-if="configuring == false"/>
      <div v-else class="authenticating">Authenticating...</div>
      <span id="new-window" tabindex="-1" class="screen-reader-text">opens in a new window</span>
   </main>

   <LibraryFooter />
   <ScrollTop />

   <Dialog v-model:visible="systemStore.showError" :modal="true" header="System Error" @hide="errorClosed()" class="error">
      {{systemStore.error}}
      <template #footer>
         <Button label="OK" autofocus class="p-button-secondary" @click="errorClosed()"/>
      </template>
   </Dialog>
</template>

<script setup>
import { onBeforeMount, ref, watch } from 'vue'
import UvaLibraryLogo from "@/components/UvaLibraryLogo.vue"
import LibraryFooter from "@/components/LibraryFooter.vue"
import { RouterView } from 'vue-router'
import Dialog from 'primevue/dialog'
import Toast from 'primevue/toast'
import ScrollTop from 'primevue/scrolltop'
import { useToast } from "primevue/usetoast"
import { useSystemStore } from "@/stores/system"

const toast = useToast()
const systemStore = useSystemStore()

const configuring = ref(true)

watch(() => systemStore.toast.show, (newShow) => {
   if ( newShow == true) {
      if ( systemStore.toast.error) {
         toast.add({severity:'error', summary:  systemStore.toast.summary, detail:  systemStore.toast.message})
      } else {
         toast.add({severity:'success', summary:  systemStore.toast.summary, detail:  systemStore.toast.message, life: 5000})
      }
      systemStore.clearToastMessage()
   }
})

onBeforeMount( async () => {
   await systemStore.getConfig()
   configuring.value = false
})

const errorClosed = (() => {
   systemStore.setError("")
   systemStore.showError = false
})
</script>

<style lang="scss">
html,
body {
   margin: 0;
   padding: 0;
   font-family: "franklin-gothic-urw", arial, sans-serif;
   -webkit-font-smoothing: antialiased;
   -moz-osx-font-smoothing: grayscale;
   color: $uva-grey-A;
   background: $uva-blue-alt-B;
   .screen-reader-text {
      position:absolute;
      left:-10000px;
      top:auto;
      width:1px;
      height:1px;
      overflow:hidden;
      visibility: hidden;
   }
}

#app {
   font-family: "franklin-gothic-urw", arial, sans-serif;
   -webkit-font-smoothing: antialiased;
   -moz-osx-font-smoothing: grayscale;
   text-align: center;
   color: $uva-text-color-base;
   margin: 0;
   padding: 0;
   background: #fafafa;
   outline: 0;
   border: 0;
}

 h1 {
   padding: 1.75rem 0;
   position: relative;
   font-weight: 700;
   color: $uva-brand-blue;
   margin: 0 ;
   font-size: 1.5em;
}
h2 {
   color: $uva-brand-blue;
   text-align: left;
   font-size: 1.3em;
   margin: 50px 0 17px 0;
}

a {
   color: $uva-brand-blue-100;
   font-weight: 500;
   text-decoration: none;
   &:hover {
      text-decoration: underline;
      color: $uva-brand-blue-200;
   }
}

a:focus, input:focus, select:focus, textarea:focus, button.pool:focus, .pre-footer a:focus  {
   outline: 2px dotted $uva-brand-blue-100;
   outline-offset: 3px;
}
a:focus {
   border-radius: 0.3rem;
}
footer, div.main-header {
   a:focus {
      outline: 2px dotted $uva-grey-200;
      outline-offset: 3px;
   }
}

header {
   background-color: $uva-brand-blue;
   color: white;
   text-align: left;
   position: relative;
   box-sizing: border-box;
   .main-header {
      display: flex;
      flex-direction: row;
      flex-wrap: nowrap;
      justify-content: space-between;
      align-content: stretch;
      align-items: center;
      gap: 10px;
   }
   p.version {
      margin: 5px 0 0 0;
      font-size: 0.5em;
      text-align: right;
   }
}
.authenticating {
   min-height:600px;
   margin-top:5%;
   font-size: 1.5rem;
   font-weight: bold;
}
 .main-header {
   display: flex;
   flex-flow: row nowrap;
   justify-content: space-between;
   align-items: center;
   div.library-link {
      width: 250px;
   }
   div.site-link {
      text-align: right;
      font-size: 1.5em;
      font-weight: bold;
      a {
         color: white;
      }
   }
}

@media only screen and (min-width: 768px) {
   .main-header {
      padding: 1vw 20px 5px 10px;
      div.site-link {
         font-size: 1.5em;
      }
   }
}
@media only screen and (max-width: 768px) {
   .main-header {
      padding: 1vw 2vw;
      div.site-link {
         font-size: 1em;
      }
   }
}
</style>
