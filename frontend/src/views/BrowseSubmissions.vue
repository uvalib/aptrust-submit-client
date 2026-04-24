<template>
   <div class="dashboard">
      <h1>Browse APTrust Submissions</h1>   
      <div class="content">
         <DataTable :value="submissionStore.searchHits" ref="hitstable" dataKey="id"
            stripedRows showGridlines responsiveLayout="scroll"
            :lazy="true" :paginator="true" :alwaysShowPaginator="false"
            @page="onPage($event)"  paginatorPosition="both"
            :first="submissionStore.offset" :rows="submissionStore.pageSize" :totalRecords="submissionStore.total"
            paginatorTemplate="PrevPageLink CurrentPageReport NextPageLink"
            currentPageReportTemplate="{first} - {last} of {totalRecords}"
            :loading="submissionStore.working"
         >
            <Column field="identifier" header="Identifier" />
            <Column field="collectionName" header="Name">
               <template #body="slotProps">
                  <span v-if="slotProps.data.collectionName">{{ slotProps.data.collectionName }}</span>
                  <span v-else class="none">Undefined</span> 
               </template>
            </Column>
            <Column field="client" header="Client" />
            <Column field="storage" header="Storage" />
            <Column field="status" header="Status">
               <template #body="slotProps">
                  <span v-if="slotProps.data.status">{{ slotProps.data.status }}</span>
                  <span v-else class="none">Unknown</span>   
               </template>
            </Column>
            <Column field="approvalEmail" header="Approval Email" class="nowrap">
               <template #body="slotProps">
                  <span v-if="slotProps.data.approvalEmail">{{ slotProps.data.approvalEmail }}</span>
                  <span v-else class="none">N/A</span>
               </template>
            </Column>
            <Column field="createdAt" header="Created" class="nowrap">
               <template #body="slotProps">{{ $formatDateTime(slotProps.data.createdAt) }}</template>
            </Column>
         </DataTable>
      </div>
   </div>
</template>

<script setup>
import { onMounted } from 'vue'
import { useSubmissionsStore } from "@/stores/submissions"
import DataTable from 'primevue/datatable';
import Column from 'primevue/column';

const submissionStore = useSubmissionsStore()

onMounted( () => {
   submissionStore.getSubmissions()
})

const onPage = ((event) => {
   submissionStore.offset = event.first
   submissionStore.getSubmissions()
})

</script>

<style lang="scss">
.dashboard {
   margin: 0 auto 50px;
   min-height: 600px;
   text-align: left;
   .none {
      color: $uva-grey-A;
      font-style: italic;
   }
   .nowrap {
      white-space: nowrap;
   }
}
@media only screen and (min-width: 768px) {
   .dashboard {
       width: 90%;
   }
}
@media only screen and (max-width: 768px) {
   .dashboard {
      width: 95%;
   }
}
</style>