<script>
  import { loading } from "../../stores/utils.js";
  import { DataTable } from "smelte";
  import { lasTxs } from "../../stores/txs.js";
</script>

<div class="rwrap flx">
  <DataTable
    {lasTxs}
    {loading}
    on:update={({ detail }) => {
      const { column, item, value } = detail;

      const index = lasTxs.findIndex(i => i.id === item.id);

      lasTxs[index][column.field] = value;
    }}
    columns={[
      { label: "ID", field: "id", class: "md:w-10", },
      {
        label: "Ep.",
        value: (v) => `S${v.season}E${v.number}`,
        class: "md:w-10",
        editable: false,
      },
      { field: "name", class: "md:w-10" },
      {
        field: "summary",
        textarea: true,
        value: v => v && v.summary ? v.summary : "",
        class: "text-sm text-gray-700 caption md:w-full sm:w-64"
      },
      {
        field: "thumbnail",
        value: (v) => v && v.image
          ? `<img src="${v.image.medium.replace("http", "https")}" height="70" alt="${v.name}">`
          : "",
        class: "w-48",
        sortable: false,
        editable: false,
      }
    ]}
  />
</div>