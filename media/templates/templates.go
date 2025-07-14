package templates

var (
	RegisterTemplate = `
		<form id="registerForm">
			<input type="text" name="username" placeholder="Username" required>
			<input type="password" name="password" placeholder="Password" required>
			<input type="text" name="email" placeholder="Email" required>
			<button type="button" onclick="submitForm()">Register</button>
		</form>
		<a href="/?action=login">Login</a>
		<a href="/catalog?page=1">Catalog</a>

		<script>
			function submitForm() {
				const form = document.getElementById("registerForm");
				const data = {
					username: form.username.value,
					password: form.password.value,
					email: form.email.value
				};

				fetch("/register", {
					method: "POST",
					headers: {
						"Content-Type": "application/json"
					},
					body: JSON.stringify(data)
				})
				.then(response => {
					if (!response.ok) {
						return response.json().then(err => { throw err; });
					}
					return response.json();
				})
				.then(data => {
					alert(data.message);
					window.location.href = data.redirect;
				})
				.catch(error => {
					alert(error.message || "Register failed");
				});
			}
		</script>
	`
	LoginTemplate = `
		<form id="loginForm">
			<input type="text" name="email" placeholder="Email" required>
			<input type="password" name="password" placeholder="Password" required>
			<button type="button" onclick="submitLoginForm()">Login</button>
		</form>
		<a href="/?action=register">Register</a>
		<a href="/catalog?page=1">Catalog</a>

		<script>
			function submitLoginForm() {
				const form = document.getElementById("loginForm");
				const data = {
					email: 	  form.email.value,
					password: form.password.value
				};

				fetch("/login", {
					method: "POST",
					headers: {
						"Content-Type": "application/json"
					},
					body: JSON.stringify(data)
				})
				.then(response => {
					if (!response.ok) {
						return response.json().then(err => { throw err; });
					}
					return response.json();
				})
				.then(data => {
					alert(data.message);
					window.location.href = data.redirect;
				})
				.catch(error => {
					alert(error.message || "Login failed");
				});
			}
		</script>`
)
