<template>
   <div class="dashboard">
      <h1>Browse APTrust Submissions</h1>   
      <div class="content">
         <div>{{  filters  }}</div>
         <DataTable :value="submissionStore.searchHits" ref="hitstable" dataKey="id"
            v-model:filters="filters" filterDisplay="menu" @filter="onFilter($event)"
            stripedRows showGridlines responsiveLayout="scroll"
            :lazy="true" :paginator="true" :alwaysShowPaginator="true"
            @page="onPage($event)"  paginatorPosition="both"
            :first="submissionStore.offset" :rows="submissionStore.pageSize" :totalRecords="submissionStore.total"
            paginatorTemplate="PrevPageLink CurrentPageReport NextPageLink"
            currentPageReportTemplate="{first} - {last} of {totalRecords}"
            :loading="submissionStore.working"
         >
            <template #header>
               <div class="search-controls">
                  <div class="search">
                     <Checkbox v-model="submissionStore.includeAutoApproved" inputId="autoapproved" binary @update:model-value="submissionStore.getSubmissions()"/>
                     <label for="autoapproved">Include auto-approved submissions</label>
                  </div>
                  <div class="search">
                     <IconField iconPosition="left" class="query">
                        <InputIcon class="pi pi-search" />
                        <InputText v-model="submissionStore.query" @keypress="searchKeyPressed($event)" fluid aria-label="search submissions"/>
                     </IconField>
                     <Button severity="secondary" label="Reset" @click="submissionStore.resetSearch()"/>
                  </div>
               </div>
            </template>
            <Column field="identifier" header="Identifier" />
            <Column field="collectionName" header="Name">
               <template #body="slotProps">
                  <span v-if="slotProps.data.collectionName">{{ slotProps.data.collectionName }}</span>
                  <span v-else class="none">Undefined</span> 
               </template>
            </Column>
            <Column field="client" header="Client" filterField="client" :showFilterMatchModes="false">
               <template #filter="{ filterModel, filterCallback }">
                  <Select v-model="filterModel.value" @change="filterCallback()" :options="clients" placeholder="Select client" />
               </template>
            </Column>
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
import { onMounted, ref, computed } from 'vue'
import { useSubmissionsStore } from "@/stores/submissions"
import { useSystemStore } from "@/stores/system"
import DataTable from 'primevue/datatable'
import Column from 'primevue/column'
import IconField from 'primevue/iconfield'
import InputIcon from 'primevue/inputicon'
import InputText from 'primevue/inputtext'
import Checkbox from 'primevue/checkbox'
import Select from 'primevue/select'
import { FilterMatchMode } from '@primevue/core/api'

const submissionStore = useSubmissionsStore()
const system = useSystemStore()

const filters = ref({
    client: { value: null, matchMode: FilterMatchMode.EQUALS },
});

const clients = computed( () => {
   let out = [] 
   system.clients.forEach( c=> {
      out.push(c.name)
   })
   return out
})


onMounted( () => {
   submissionStore.getSubmissions()
})

function onFilter(event) {
   submissionStore.offset = 0
   submissionStore.filters = []
   Object.entries(event.filters).forEach(([key, data]) => {
      if (data.value && data.value != "") {
         submissionStore.filters.push({field: key, match: data.matchMode, value: data.value})
      }
   })
   submissionStore.getSubmissions()
}

const onPage = ((event) => {
   submissionStore.offset = event.first
   submissionStore.getSubmissions()
})

const searchKeyPressed = ((event) => {
   submissionStore.offset = 0
   if (event.keyCode == 13) {
      submissionStore.getSubmissions()
   }
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
   .search-controls {
      display: flex;
      flex-flow: row nowrap;
      justify-content: space-between;
      .search {
         display: flex;
         flex-flow: row nowrap;   
         align-items: center;
         justify-content: flex-start;
         gap: 10px;
      }
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