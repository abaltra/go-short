<script>
	export let origin_url = "";
	export let dest_url = ""; 

	async function minimize() {
		console.log(`Minimizing ${origin_url}`)

		let response = await fetch(`http://localhost:8000`, {
			method:'POST',
			body: JSON.stringify({
				URL: origin_url
			})
		})

		let body = await response.json()

		console.log(body)
		dest_url = `http://localhost:8000/${body.Code}`
	}

	function forward() {
		location.assign(dest_url)
	}
	
</script>

<main>
	<h1>Hello there!</h1>
	<p>Let's minimize some URLs</p>
	<input bind:value="{origin_url}"/>
	<button on:click="{minimize}">Minimize!</button>
	<p></p>
	<input bind:value="{dest_url}"/>
	<button on:click="{forward}">Forward!</button>
</main>

<style>
	main {
		text-align: center;
		padding: 1em;
		max-width: 240px;
		margin: 0 auto;
	}

	h1 {
		color: #ff3e00;
		text-transform: uppercase;
		font-size: 4em;
		font-weight: 100;
	}

	@media (min-width: 640px) {
		main {
			max-width: none;
		}
	}
</style>