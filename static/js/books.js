// 数据存储示例
let books = [
    { ISBN: "123456", Title: "Book 1", Author: "Author 1", Publisher: "Publisher 1", Description: "Description 1" },
    { ISBN: "654321", Title: "Book 2", Author: "Author 2", Publisher: "Publisher 2", Description: "Description 2" },
];

// 动态加载书籍到表格
function loadBooks() {
    const bookList = document.getElementById("bookList");
    bookList.innerHTML = books.map(book => `
        <tr>
            <td><input type="checkbox" value="${book.ISBN}"></td>
            <td>${book.ISBN}</td>
            <td>${book.Title}</td>
            <td>${book.Author}</td>
            <td>${book.Publisher}</td>
            <td>${book.Description}</td>
        </tr>
    `).join('');
}

// 添加书籍功能
function showAddBookForm() {
    const isbn = prompt("Enter ISBN:");
    const title = prompt("Enter Title:");
    const author = prompt("Enter Author:");
    const publisher = prompt("Enter Publisher:");
    const description = prompt("Enter Description:");

    if (isbn && title && author && publisher) {
        books.push({ ISBN: isbn, Title: title, Author: author, Publisher: publisher, Description: description });
        loadBooks();
        alert("Book added successfully!");
    } else {
        alert("All fields are required!");
    }
}

// 编辑书籍功能
function showEditBookForm() {
    const selectedBook = document.querySelector('input[type="checkbox"]:checked');
    if (!selectedBook) {
        alert("Please select a book to edit.");
        return;
    }

    const book = books.find(b => b.ISBN === selectedBook.value);
    if (book) {
        book.Title = prompt("Edit Title:", book.Title) || book.Title;
        book.Author = prompt("Edit Author:", book.Author) || book.Author;
        book.Publisher = prompt("Edit Publisher:", book.Publisher) || book.Publisher;
        book.Description = prompt("Edit Description:", book.Description) || book.Description;
        loadBooks();
        alert("Book updated successfully!");
    }
}

// 删除书籍功能
function deleteSelectedBooks() {
    const selectedBooks = Array.from(document.querySelectorAll('input[type="checkbox"]:checked'));
    if (selectedBooks.length === 0) {
        alert("Please select at least one book to delete.");
        return;
    }

    const selectedISBNs = selectedBooks.map(book => book.value);
    books = books.filter(book => !selectedISBNs.includes(book.ISBN));
    loadBooks();
    alert("Selected books deleted successfully!");
}

// 搜索书籍功能
function searchBooks() {
    const keyword = document.getElementById("searchKeyword").value.toLowerCase();
    const filteredBooks = books.filter(book =>
        book.Title.toLowerCase().includes(keyword) || book.Description.toLowerCase().includes(keyword)
    );

    const bookList = document.getElementById("bookList");
    bookList.innerHTML = filteredBooks.map(book => `
        <tr>
            <td><input type="checkbox" value="${book.ISBN}"></td>
            <td>${book.ISBN}</td>
            <td>${book.Title}</td>
            <td>${book.Author}</td>
            <td>${book.Publisher}</td>
            <td>${book.Description}</td>
        </tr>
    `).join('');
}

// 页面加载完成时初始化书籍列表
document.addEventListener("DOMContentLoaded", function () {
    loadBooks();
    document.getElementById("addButton").addEventListener("click", showAddBookForm);
    document.getElementById("editButton").addEventListener("click", showEditBookForm);
    document.getElementById("deleteButton").addEventListener("click", deleteSelectedBooks);
    document.getElementById("searchButton").addEventListener("click", searchBooks);
});