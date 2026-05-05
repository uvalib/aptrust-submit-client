export default {
   install: (app) => {
      app.config.globalProperties.$formatFileSize = ( bytes ) => {
         let sz =  bytes / 1000.0 / 1000.0
         let units = "MB"
         if ( sz > 1000 ) {
            sz = sz / 1000.0
            units = "GB"
         }
         return `${ sz.toFixed(2)} ${units}`
      }
   }
}