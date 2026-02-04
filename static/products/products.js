const productsContainer = document.getElementById('products');
const urlParams = new URLSearchParams(window.location.search);
const query = urlParams.get('q') || '';

let filters = { q: query };
let currentPage = 1;
let isLoading = false;
let hasMore = true;

function setupInfiniteScroll() {
    window.addEventListener('scroll', () => {
        const { scrollTop, scrollHeight, clientHeight } = document.documentElement;
        if (scrollTop + clientHeight >= scrollHeight - 100 && !isLoading && hasMore) {
            loadMoreProducts();
        }
    });
}

async function loadProducts(clear = true) {
    if (clear) {
        currentPage = 1;
        hasMore = true;
        if (productsContainer) productsContainer.innerHTML = '';
    }

    const token = localStorage.getItem('token');
    const params = new URLSearchParams();
    if (filters.q) params.append('q', filters.q);
    if (filters.brand) params.append('brand', filters.brand);
    if (filters.availability) params.append('availability', filters.availability);
    if (filters.sort_price) params.append('sort_price', filters.sort_price);
    params.append('page', currentPage);
    params.append('limit', '20');

    isLoading = true;

    try {
        const res = await fetch(`/api/products?${params}`, {
            headers: token ? { 'Authorization': 'Bearer ' + token } : {}
        });

        if (!res.ok) {
            if (res.status === 401) localStorage.removeItem('token');
            if (productsContainer && currentPage === 1) productsContainer.innerHTML = '<p class="error">Ошибка загрузки товаров</p>';
            hasMore = false;
            return;
        }

        const products = await res.json();

        if (!products || products.length === 0) {
            if (currentPage === 1 && productsContainer) productsContainer.innerHTML = '<p>Товары не найдены</p>';
            hasMore = false;
            return;
        }

        if (products.length < 20) {
            hasMore = false;
        }
        showProducts(products, clear);
        currentPage += 1;
    } catch (error) {
        console.error('Ошибка загрузки товаров:', error);
        if (productsContainer && currentPage === 1) productsContainer.innerHTML = '<p>Ошибка соединения с сервером</p>';
    } finally {
        isLoading = false;
    }
}

async function loadMoreProducts() {
    if (isLoading || !hasMore) return;
    await loadProducts(false);
}

function showProducts(products, clear = true) {
    if (!productsContainer) return;
    if (clear) productsContainer.innerHTML = '';

    products.forEach(product => {
        const photo = product.photos_urls && product.photos_urls[0] ? product.photos_urls[0] : '../assets/placeholder.jpg';

        const productCard = document.createElement('div');
        productCard.className = 'product-card';
        productCard.dataset.id = product.id;

        productCard.innerHTML = `
            <div class="product-image-container">
                <img src="${photo}" alt="${product.name}">
                <div class="favorite-icon" data-id="${product.id}">
                    <svg width="38" height="38" viewBox="0 0 38 38" fill="none">
                        <path d="M32.9967 7.29917C32.188 6.4901 31.2278 5.84828 30.171 5.41039C29.1142 4.9725 27.9814 4.74712 26.8375 4.74712C25.6936 4.74712 24.5608 4.9725 23.504 5.41039C22.4472 5.84828 21.487 6.4901 20.6783 7.29917L19 8.97751L17.3217 7.29917C15.6882 5.66566 13.4726 4.74796 11.1625 4.74796C8.85237 4.74796 6.63685 5.66566 5.00334 7.29917C3.36982 8.93269 2.45213 11.1482 2.45213 13.4583C2.45213 15.7685 3.36982 17.984 5.00334 19.6175L19 33.6142L32.9967 19.6175C33.8057 18.8088 34.4476 17.8486 34.8854 16.7918C35.3233 15.735 35.5487 14.6023 35.5487 13.4583C35.5487 12.3144 35.3233 11.1817 34.8854 10.1249C34.4476 9.06805 33.8057 8.10787 32.9967 7.29917Z" stroke="#5A5A6E" stroke-width="4"/>
                    </svg>
                </div>
            </div>
            <div class="card-body">
                <h5>${product.name}</h5>
                ${product.brand ? `<p class="brand">${product.brand}</p>` : ''}
                <p class="price">${product.price} ₽</p>
                ${product.availability !== undefined ? 
                    `<p class="availability ${product.availability ? 'in-stock' : 'out-of-stock'}">
                        ${product.availability ? 'В наличии' : 'Нет в наличии'}
                    </p>` : ''
                }
                ${product.description ? `<p class="description">${product.description}</p>` : ''}
            </div>
        `;

        productsContainer.appendChild(productCard);
    });

    addClickHandlers();
}

function addClickHandlers() {
    document.querySelectorAll('.favorite-icon').forEach(icon => {
        icon.addEventListener('click', async function() {
            const productId = this.dataset.id;
            const token = localStorage.getItem('token');

            if (!token) {
                const modalElement = document.getElementById('authModal');
                if (modalElement) {
                    const modal = new bootstrap.Modal(modalElement);
                    modal.show();
                }
                return;
            }

            try {
                const res = await fetch(`/api/like?product_id=${productId}`, {
                    method: 'POST',
                    headers: { 'Authorization': 'Bearer ' + token }
                });

                if (!res.ok) {
                    if (res.status === 401) {
                        localStorage.removeItem('token');
                        const modalElement = document.getElementById('authModal');
                        if (modalElement) {
                            const modal = new bootstrap.Modal(modalElement);
                            modal.show();
                        }
                    }
                    return;
                }

                const data = await res.json();
                const heart = this.querySelector('path');

                if (data.status === 'лайк добавлен') {
                    heart.setAttribute('fill', '#5A87E6');
                    heart.setAttribute('stroke', '#5A87E6');
                } else if (data.status === 'лайк удалён') {
                    heart.removeAttribute('fill');
                    heart.setAttribute('stroke', '#5A5A6E');
                }

            } catch (error) {
                console.error('Ошибка при добавлении в избранное:', error);
            }
        });
    });
}

async function loadAndPopulateBrands() {
    try {
        const token = localStorage.getItem('token');
        const response = await fetch('/api/products?limit=100', {
            headers: token ? { 'Authorization': 'Bearer ' + token } : {}
        });

        if (response.ok) {
            const products = await response.json();
            const brandsSet = new Set();
            products.forEach(product => {
                if (product.brand) brandsSet.add(product.brand);
            });

            const brands = Array.from(brandsSet).sort();
            const brandSelect = document.querySelector('#brandFilter');
            if (brandSelect) {
                brandSelect.innerHTML = '<option value="">Все бренды</option>';
                brands.forEach(brand => {
                    const option = document.createElement('option');
                    option.value = brand;
                    option.textContent = brand;
                    brandSelect.appendChild(option);
                });
            }
        }
    } catch (error) {
        console.error('Ошибка загрузки брендов:', error);
    }
}

function applyAllFilters() {
    const searchValue = document.querySelector('#searchFilter')?.value.trim();
    filters.q = searchValue || '';

    const brandSelect = document.querySelector('#brandFilter');
    filters.brand = brandSelect && brandSelect.value ? brandSelect.value : '';

    const availabilitySelect = document.querySelector('#availabilityFilter');
    filters.availability = availabilitySelect && availabilitySelect.value ? availabilitySelect.value : '';

    const sortSelect = document.querySelector('#sortFilter');
    filters.sort_price = sortSelect && sortSelect.value ? sortSelect.value : '';

    loadProducts();
}

function resetAllFilters() {
    filters = { q: query };

    const searchInput = document.querySelector('#searchFilter');
    if (searchInput) searchInput.value = query || '';

    const brandSelect = document.querySelector('#brandFilter');
    if (brandSelect) brandSelect.value = '';

    const availabilitySelect = document.querySelector('#availabilityFilter');
    if (availabilitySelect) availabilitySelect.value = '';

    const sortSelect = document.querySelector('#sortFilter');
    if (sortSelect) sortSelect.value = '';

    loadProducts();
}

setupInfiniteScroll();
loadAndPopulateBrands();
loadProducts();