document.addEventListener('DOMContentLoaded', () => {
    const medicationsList = document.getElementById('medicationsList');
    const searchInput = document.getElementById('search');
    const addBtn = document.getElementById('addBtn');
    const authBtn = document.getElementById('authBtn');
    const modal = document.getElementById('modal');
    const adjustModal = document.getElementById('adjustModal');
    const authModal = document.getElementById('authModal');
    const closeBtns = document.querySelectorAll('.close');
    const medicationForm = document.getElementById('medicationForm');
    const adjustForm = document.getElementById('adjustForm');
    const authForm = document.getElementById('authForm');
    const prevPage = document.getElementById('prevPage');
    const nextPage = document.getElementById('nextPage');
    const pageIndicator = document.getElementById('pageIndicator');

    let currentPage = 0;
    const itemsPerPage = 5;
    let medications = [];
    let adjustMedicationId = null;
    let adjustAction = '';

    let authData = null; // Данные авторизации

    // Добавляем обработчик для кнопки авторизации
    authBtn.addEventListener('click', () => {
        authModal.style.display = 'block';
    });

    // Обработчик для формы авторизации
    authForm.addEventListener('submit', (event) => {
        event.preventDefault();

        const formData = new FormData(authForm);
        const username = formData.get('username');
        const password = formData.get('password');

        // Сохраняем данные авторизации на фронтенде
        authData = {
            username: username,
            password: password
        };

        // Сбрасываем форму и закрываем модальное окно
        authModal.style.display = 'none';
        authForm.reset();

        // Перезагружаем данные после авторизации
        fetchMedications();
    });

    // Функция для отправки запросов с авторизацией
    function fetchWithAuth(url, options) {
        // Если есть данные авторизации, добавляем Basic Auth заголовок
        if (authData) {
            const headers = options.headers || {};
            headers['Authorization'] = `Basic ${btoa(`${authData.username}:${authData.password}`)}`;
            options.headers = headers;
        }

        // Отправляем запрос с обновленными опциями
        return fetch(url, options);
    }

    // Функция для получения списка медикаментов
    function fetchMedications(page = currentPage, limit = itemsPerPage, searchQuery = '') {
        var url = '';
        if (searchQuery === '') {
            url = '/api/getPage';
            const requestData = {
                page: page,
                limit: limit,
                searchQuery: searchQuery
            };

        fetchWithAuth(url, {
            method: 'POST',
            headers: {
                'Content-Type': 'application/json'
            },
            body: JSON.stringify(requestData)
        })
        .then(response => response.json())
        .then(data => {
            medications = data;
            if (!medications && currentPage > 0) {
                currentPage--;
                fetchMedications(currentPage, itemsPerPage, searchQuery);
            } else {
                renderMedications();
                updatePaginationControls();
            }
        })
        .catch(error => console.error('Error fetching medications:', error));
        }else{
            url = '/api/searchDrug';
            const requestData = {
                name: searchQuery,
            };

        fetchWithAuth(url, {
            method: 'POST',
            headers: {
                'Content-Type': 'application/json'
            },
            body: JSON.stringify(requestData)
        })
        .then(response => response.json())
        .then(data => {
            medications = data;
            if (!medications && currentPage > 0) {
                currentPage--;
                fetchMedications(currentPage, itemsPerPage, searchQuery);
            } else {
                renderMedications();
                updatePaginationControls();
            }
        })
        .catch(error => console.error('Error fetching medications:', error));
        }

        
    }

    // Функция для отображения списка медикаментов в таблице
    function renderMedications() {
        medicationsList.innerHTML = '';
        medications.forEach(medication => {
            const row = document.createElement('tr');
            row.innerHTML = `
                <td>${medication.name}</td>
                <td>${medication.quantity}</td>
                <td>
                    <button onclick="deleteMedication('${medication.name}')">Удалить</button>
                    <button onclick="showAdjustModal('${medication.name}', 'add')">Добавить</button>
                    <button onclick="showAdjustModal('${medication.name}', 'sub')">Забрать</button>
                </td>
            `;
            medicationsList.appendChild(row);
        });
        pageIndicator.textContent = `Страница ${currentPage + 1}`;
    }

    // Функция для обновления элементов управления пагинацией
    function updatePaginationControls() {
        prevPage.disabled = currentPage === 0;
        nextPage.disabled = medications.length < itemsPerPage;
    }

    // Обработчик для поиска при вводе в поле поиска
    searchInput.addEventListener('keydown', (event) => {
        if (event.key === 'Enter') {
            event.preventDefault();
            const searchQuery = searchInput.value;
            currentPage = 0;  // Сброс страницы на первую при новом поиске
            fetchMedications(currentPage, itemsPerPage, searchQuery);
        }
    });

    // Обработчик для открытия модального окна добавления медикамента
    addBtn.addEventListener('click', () => {
        modal.style.display = 'block';
        medicationForm.reset();
    });

    // Обработчики для закрытия всех модальных окон
    closeBtns.forEach(btn => {
        btn.addEventListener('click', () => {
            modal.style.display = 'none';
            adjustModal.style.display = 'none';
            authModal.style.display = 'none';
        });
    });

    // Обработчик для отправки формы добавления медикамента
    medicationForm.addEventListener('submit', (event) => {
        event.preventDefault();

        const formData = new FormData(medicationForm);
        const medication = {
            name: formData.get('name'),
            quantity: Number(formData.get('quantity')),
        };

        fetchWithAuth('/api/newDrug', {
            method: 'POST',
            headers: {
                'Content-Type': 'application/json'
            },
            body: JSON.stringify(medication)
        })
        .then(() => {
            modal.style.display = 'none';
            fetchMedications(currentPage, itemsPerPage);
        })
        .catch(error => console.error('Error adding medication:', error));
    });

    // Обработчик для кнопки "Предыдущая страница"
    prevPage.addEventListener('click', () => {
        if (currentPage > 0) {
            currentPage--;
            fetchMedications(currentPage, itemsPerPage, searchInput.value);
        }
    });

    // Обработчик для кнопки "Следующая страница"
    nextPage.addEventListener('click', () => {
        if (medications.length === itemsPerPage) {
            currentPage++;
            fetchMedications(currentPage, itemsPerPage, searchInput.value);
        }
    });

    // Функция для удаления медикамента
    window.deleteMedication = (name) => {
        if (confirm('Вы уверены, что хотите удалить этот медикамент?')) {
            fetchWithAuth('/api/deleteDrug', {
                method: 'POST',
                headers: {
                    'Content-Type': 'application/json'
                },
                body: JSON.stringify({ name: name })
            })
            .then(() => {
                fetchMedications(currentPage, itemsPerPage, searchInput.value);
            })
            .catch(error => console.error('Error deleting medication:', error));
        }
    };

    // Функция для отображения модального окна для корректировки количества медикамента
    window.showAdjustModal = (name, action) => {
        adjustMedicationId = name;
        adjustAction = action;
        adjustModal.style.display = 'block';
        adjustForm.reset();
    };

    // Обработчик для отправки формы корректировки количества медикамента
    adjustForm.addEventListener('submit', (event) => {
    event.preventDefault();
    const adjustQuantity = Number(document.getElementById('adjustQuantity').value);
    const medication = {
        name: adjustMedicationId,
        quantity: Number(adjustQuantity),
    };

    const endpoint = adjustAction === 'sub' ? '/api/subDrug' : '/api/addDrug';

    fetchWithAuth(endpoint, {
        method: 'POST',
        headers: {
            'Content-Type': 'application/json'
        },
        body: JSON.stringify(medication)
    })
    .then(() => {
        adjustModal.style.display = 'none';
        fetchMedications(currentPage, itemsPerPage, searchInput.value);
    })
    .catch(error => console.error('Error adjusting medication:', error));
});
});
