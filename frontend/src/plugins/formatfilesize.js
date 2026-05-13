export default {
   install: (app) => {
      app.config.globalProperties.$formatFileSize = ( bytes ) => {
         let units = ["KB", "MB", "GB", "TB"]
         let sz = bytes / 1000.0
         let idx = 0
         while ( sz > 1000) {
            idx++
            sz = sz / 1000.0
         }
         return `${ sz.toFixed(2)} ${units[idx]}`
      }
   }
}