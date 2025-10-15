// Utilities
const api = {
    items: "/api/items",
    analytics: "/api/analytics",
};

const qs = (sel, root = document) => root.querySelector(sel);
const qsa = (sel, root = document) => Array.from(root.querySelectorAll(sel));

function toJSON(form) {
    const data = new FormData(form);
    return Object.fromEntries([...data.entries()].map(([k, v]) => [k, v instanceof File ? v : String(v).trim()]));
}

function formatMoney(n) {
    if (n === null || n === undefined || Number.isNaN(Number(n))) return "—";
    return Number(n).toLocaleString("ru-RU", { minimumFractionDigits: 2, maximumFractionDigits: 2 });
}

function downloadCSV(filename, rows) {
    const escapeCSV = (value) => {
        const s = String(value ?? "");
        if (/[",;\n]/.test(s)) return '"' + s.replace(/"/g, '""') + '"';
        return s;
    };
    const csv = rows.map(r => r.map(escapeCSV).join(",")).join("\n");
    const blob = new Blob(["\uFEFF" + csv], { type: "text/csv;charset=utf-8;" });
    const url = URL.createObjectURL(blob);
    const a = document.createElement("a");
    a.href = url;
    a.download = filename;
    document.body.appendChild(a);
    a.click();
    a.remove();
    URL.revokeObjectURL(url);
}

// Items CRUD
async function fetchItems(params = {}) {
    const url = new URL(api.items, window.location.origin);
    if (params.from) url.searchParams.set("from", params.from);
    if (params.to) url.searchParams.set("to", params.to);
    if (params.sort_by) url.searchParams.append("sort_by", params.sort_by);
    const res = await fetch(url.toString());
    if (!res.ok) throw new Error(await res.text());
    return res.json();
}

async function createItem(payload) {
    const res = await fetch(api.items, {
        method: "POST",
        headers: { "Content-Type": "application/json" },
        body: JSON.stringify(payload),
    });
    if (!res.ok) throw new Error(await res.text());
    return res.json();
}

async function updateItem(id, payload) {
    const res = await fetch(`${api.items}/${id}`, {
        method: "PUT",
        headers: { "Content-Type": "application/json" },
        body: JSON.stringify(payload),
    });
    if (!res.ok) throw new Error(await res.text());
    return res.json();
}

async function deleteItem(id) {
    const res = await fetch(`${api.items}/${id}`, { method: "DELETE" });
    if (!res.ok && res.status !== 204) throw new Error(await res.text());
}

// Analytics
async function fetchAnalytics(from, to) {
    const url = new URL(api.analytics, window.location.origin);
    url.searchParams.set("from", from);
    url.searchParams.set("to", to);
    const res = await fetch(url.toString());
    if (!res.ok) throw new Error(await res.text());
    return res.json();
}

// UI rendering
function renderRows(items) {
    const tbody = qs("#items-table tbody");
    tbody.innerHTML = "";
    const tmpl = qs("#row-template");
    items.forEach(item => {
        const tr = tmpl.content.firstElementChild.cloneNode(true);
        qs(".cell-id", tr).textContent = item.id;
        qs(".cell-type", tr).textContent = item.type === "income" ? "Доход" : "Расход";
        qs(".cell-amount", tr).textContent = formatMoney(item.amount);
        qs(".cell-date", tr).textContent = item.date;
        qs(".cell-category", tr).textContent = item.category;
        qs(".cell-description", tr).textContent = item.description || "";

        qs(".edit", tr).addEventListener("click", () => openEditRow(tr, item));
        qs(".delete", tr).addEventListener("click", async () => {
            if (!confirm("Удалить запись?")) return;
            await deleteItem(item.id);
            await reloadItems();
        });
        tbody.appendChild(tr);
    });
}

function openEditRow(tr, item) {
    const original = tr.cloneNode(true);
    const cells = qsa("td", tr);
    cells[1].innerHTML = `<select class="edit-type"><option value="income">Доход</option><option value="expense">Расход</option></select>`;
    cells[2].innerHTML = `<input class="edit-amount" type="number" step="0.01" min="0.01" value="${item.amount}">`;
    cells[3].innerHTML = `<input class="edit-date" type="date" value="${item.date}">`;
    cells[4].innerHTML = `<input class="edit-category" type="text" value="${item.category}">`;
    cells[5].innerHTML = `<input class="edit-description" type="text" maxlength="255" value="${item.description ?? ""}">`;
    const typeSel = qs(".edit-type", tr);
    typeSel.value = item.type;
    const actions = cells[6];
    actions.innerHTML = "";
    const btnSave = document.createElement("button");
    btnSave.className = "link";
    btnSave.textContent = "Сохранить";
    const btnCancel = document.createElement("button");
    btnCancel.className = "link";
    btnCancel.textContent = "Отмена";
    actions.append(btnSave, btnCancel);

    btnCancel.addEventListener("click", () => {
        tr.replaceWith(original);
    });

    btnSave.addEventListener("click", async () => {
        const payload = {
            type: typeSel.value,
            amount: Number(qs(".edit-amount", tr).value),
            date: qs(".edit-date", tr).value,
            category: qs(".edit-category", tr).value.trim(),
            description: qs(".edit-description", tr).value.trim(),
        };
        await updateItem(item.id, payload);
        await reloadItems();
    });
}

async function reloadItems() {
    const filters = toJSON(qs("#filters"));
    const items = await fetchItems({
        from: filters.from || undefined,
        to: filters.to || undefined,
        sort_by: filters.sort_by || undefined,
    });
    renderRows(items);
}

async function reloadAnalytics() {
    const f = qs("#analytics-filters");
    if (!f.from.value || !f.to.value) return;
    const data = await fetchAnalytics(f.from.value, f.to.value);
    qs("#kpi-sum").textContent = formatMoney(data.sum);
    qs("#kpi-avg").textContent = formatMoney(data.average);
    qs("#kpi-count").textContent = Number(data.count).toLocaleString("ru-RU");
    qs("#kpi-median").textContent = formatMoney(data.median);
    qs("#kpi-p90").textContent = formatMoney(data.percentile);
}

// Events
document.addEventListener("DOMContentLoaded", () => {
    // Today as default dates
    const today = new Date().toISOString().slice(0, 10);
    [qs('#filters [name="to"]'), qs('#analytics-filters [name="to"]')].forEach(el => { if (el) el.value = today; });
    const monthAgo = new Date();
    monthAgo.setMonth(monthAgo.getMonth() - 1);
    const monthAgoStr = monthAgo.toISOString().slice(0, 10);
    [qs('#filters [name="from"]'), qs('#analytics-filters [name="from"]')].forEach(el => { if (el) el.value = monthAgoStr; });

    // Create item
    const form = qs('#item-form');
    form.addEventListener('submit', async (e) => {
        e.preventDefault();
        const err = qs('#form-error');
        err.hidden = true;
        try {
            const payload = toJSON(form);
            payload.amount = Number(payload.amount);
            await createItem(payload);
            form.reset();
            await reloadItems();
        } catch (e) {
            err.textContent = String(e);
            err.hidden = false;
        }
    });

    // Filters
    qs('#filters').addEventListener('submit', async (e) => {
        e.preventDefault();
        await reloadItems();
    });
    qs('#clear-filters').addEventListener('click', async () => {
        const f = qs('#filters');
        f.reset();
        await reloadItems();
    });

    // CSV export items
    qs('#export-csv').addEventListener('click', async () => {
        const filters = toJSON(qs('#filters'));
        const items = await fetchItems({
            from: filters.from || undefined,
            to: filters.to || undefined,
            sort_by: filters.sort_by || undefined,
        });
        const rows = [
            ["id", "type", "amount", "date", "category", "description", "created_at", "updated_at"],
            ...items.map(i => [i.id, i.type, i.amount, i.date, i.category, i.description ?? "", i.created_at, i.updated_at]),
        ];
        downloadCSV(`items_${Date.now()}.csv`, rows);
    });

    // Analytics
    qs('#analytics-filters').addEventListener('submit', async (e) => {
        e.preventDefault();
        await reloadAnalytics();
    });
    qs('#export-analytics-csv').addEventListener('click', async () => {
        const f = qs('#analytics-filters');
        if (!f.from.value || !f.to.value) return;
        const d = await fetchAnalytics(f.from.value, f.to.value);
        const rows = [["sum", "average", "count", "median", "percentile"], [d.sum, d.average, d.count, d.median, d.percentile]];
        downloadCSV(`analytics_${f.from.value}_${f.to.value}.csv`, rows);
    });

    // Initial load
    reloadItems();
    reloadAnalytics();
});


