<template>
   <Button size="small" severity="secondary" @click="showDialog = true" label="Details" iconPos="right"/>
   <Dialog v-model:visible="showDialog" :modal="true" header="Submission Failures" >
      <DataTable :value="submission.detail.failures" dataKey="id" 
         stripedRows showGridlines responsiveLayout="scroll" :alwaysShowPaginator="false"
         :lazy="false" :paginator="true" :rows="30" :totalRecords="submission.detail.failures.length"
         paginatorTemplate="FirstPageLink PrevPageLink CurrentPageReport NextPageLink LastPageLink RowsPerPageDropdown"
         currentPageReportTemplate="{first} - {last} of {totalRecords}" paginatorPosition="both"
      >
         <Column field="createdAt" header="Timestamp">
            <template #body="slotProps">
               {{ $formatDateTime(slotProps.data.createdAt) }}
            </template>
         </Column>
         <Column field="failure" header="Failure" />
      </DataTable>       
   </Dialog>
</template>

<script setup>
import Dialog from 'primevue/dialog'
import { ref } from 'vue'
import DataTable from 'primevue/datatable'
import Column from 'primevue/column'
import { useSubmissionsStore } from "@/stores/submissions"

const submission = useSubmissionsStore()
const showDialog = ref(false)
</script>

<style lang="scss" scoped>
</style>
