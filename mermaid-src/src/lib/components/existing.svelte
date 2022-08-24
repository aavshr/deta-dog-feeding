<script lang="ts">
    import Card from '$lib/components/card/card.svelte';
    import { API } from '../../api';
    import { updateCode } from '$lib/util/state';

    const api = new API(`https://${window.location.host}/backend`)
    let codes = []

    const loadCodes = async () => {
       codes = await api.listCodes()
    }

    const loadCode = (c) => {
        updateCode(c.Content, {
            updateDiagram: true,
            updateEditor: true,
            resetPanZoom: true
        });
    };

</script>

<Card title="Your Diagrams" isOpen="{false}">
    <div class="flex gap-2 flex-wrap p-2">
        {#each codes as c}
            <button class="btn btn-primary normal-case btn-sm" on:click={() => loadCode(c)}
            >{c.Key}</button>
        {/each}
    </div>
</Card>