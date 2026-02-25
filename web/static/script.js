const apiBase = '/api/products';

async function loadProducts() {
    const response = await fetch(apiBase);
    const products = await response.json();
    const tbody = document.querySelector('#productTable tbody');
    tbody.innerHTML = '';
    products.forEach(product => {
        const row = tbody.insertRow();
        row.insertCell().textContent = product.name;
        row.insertCell().textContent = product.count;
        const dateStr = product.date.substring(0,10);
        row.insertCell().textContent = dateStr;
        const deleteCell = row.insertCell();
        const deleteBtn = document.createElement('button');
        deleteBtn.textContent = 'Удалить';
        deleteBtn.onclick = () => deleteProduct(product.id);
        deleteCell.appendChild(deleteBtn);
    });
}

async function addProduct(event) {
    event.preventDefault();
    const name = document.getElementById('name').value;
    const count = parseInt(document.getElementById('count').value);
    const date = document.getElementById('date').value;

    const response = await fetch(apiBase, {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify({ name, count, date })
    });
    if (response.ok) {
        document.getElementById('addForm').reset();
        loadProducts();
    } else {
        const err = await response.json();
        alert('Ошибка добавления: ' + (err.error || 'неизвестная ошибка'));
    }
}

async function deleteProduct(id) {
    if (!confirm('Удалить продукт?')) return;
    const response = await fetch(`${apiBase}/${id}`, { method: 'DELETE' });
    if (response.ok) {
        loadProducts();
    } else {
        const err = await response.json();
        alert('Ошибка удаления: ' + (err.error || 'неизвестная ошибка'));
    }
}

document.getElementById('addForm').addEventListener('submit', addProduct);
loadProducts();