<template>
	<div>
		<!-- Drawer -->
		<v-navigation-drawer clipped app :mini-variant="drawer">
			<v-list>
				<v-list-tile v-for="item in items" :to="item.to" :key="item.id" active-class="white--text">
					<v-list-tile-action>
						<v-icon medium>{{item.icon}}</v-icon>
					</v-list-tile-action>
					<v-list-tile-content>
						<v-list-tile-title>{{item.title}}</v-list-tile-title>
					</v-list-tile-content>
				</v-list-tile>
			</v-list>
		</v-navigation-drawer>
		<!-- Toolbar -->
		<v-toolbar app fixed clipped-left>
			<v-toolbar-side-icon @click.stop="drawer = !drawer"></v-toolbar-side-icon>
			<img src="/images/logo.png" class="fw-logo">
			<v-toolbar-title>FireWatch</v-toolbar-title>
			<v-spacer></v-spacer>
			<v-toolbar-items class="hidden-sm-down">
				<v-menu offset-y>
					<template v-slot:activator="{ on }">
						<v-btn dark v-on="on">
							{{username}}
							<v-icon>keyboard_arrow_down</v-icon>
						</v-btn>
					</template>
					<v-list>
						<v-list-tile
							v-for="(item, index) in accountItems"
							:key="index"
							:to="item.to"
							active-class="white--text"
						>
							<v-list-tile-title>
								<v-icon>{{item.icon}}</v-icon>
								{{ item.title }}
							</v-list-tile-title>
						</v-list-tile>
					</v-list>
				</v-menu>
			</v-toolbar-items>
		</v-toolbar>

		<v-content>
			<v-container fluid>
				<router-view name="int"/>
			</v-container>
		</v-content>
	</div>
</template>

<script>
	export default {
		name: "Navigation",
		data: () => ({
			drawer: null,
			items: [
				{ title: "Dashboard", icon: "dashboard", to: "/" },
				{ title: "Devices", icon: "memory", to: "/device/" },
				{ title: "Settings", icon: "settings", to: "/settings/" }
			],
			accountItems: [
				{ title: "Account", icon: "account_circle", to: "/account/" },
				{ title: "Log Out", icon: "power_settings_new", to: "/logout/" }
			],
			username: JSON.parse(localStorage.getItem("user")).UserName
		})
	};
</script>

<style scoped>
	.fw-logo {
		height: 2em;
		width: auto;
		margin-left: 1em;
		margin-right: -0.5em;
	}
</style>
