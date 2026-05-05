import { createApp, markRaw } from 'vue'
import { createPinia } from 'pinia'
import formatDateTimePlugin from './plugins/formatdatetime'
import formatFileSizePlugin from './plugins/formatfilesize'

import App from './App.vue'
import router from './router'

const app = createApp(App)

app.use(router)
app.use(formatDateTimePlugin)
app.use(formatFileSizePlugin)

// Primevue setup
import PrimeVue from 'primevue/config'
import UVA from './assets/theme/uva'
import ConfirmationService from 'primevue/confirmationservice'
import ToastService from 'primevue/toastservice'
import Button from 'primevue/button'
import ConfirmDialog from 'primevue/confirmdialog'
import 'primeicons/primeicons.css'

app.use(PrimeVue, {
   theme: {
      preset: UVA,
      options: {
         prefix: 'p',
         darkModeSelector: '.aptsubmit-dark'
      }
   }
})

app.use(ConfirmationService)
app.use(ToastService)

app.component("Button", Button)
app.component("ConfirmDialog", ConfirmDialog)


// Per some suggestions on vue / pinia git hub issue reports, create and add pinia support LAST
// and use the chained form of the setup. This to avid problems where the vuew dev tools fail to
// include pinia in the tools
app.use(createPinia().use( ({ store }) => {
   store.router = markRaw(router)
}))

app.mount('#app')
