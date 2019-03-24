<template>
	<v-container>
		<v-layout row wrap>
			<v-flex xs12 md8>
				<h1>User Details</h1>
			</v-flex>
			<v-flex xs12 md4>
				<h1>Account Actions</h1>
				<v-layout column>
					<v-btn color="info" v-on:click.native="LogOutAllDevices">
						<v-icon>power_settings_new</v-icon>Log Out from all Devices
					</v-btn>

					<v-dialog v-model="dialog" width="500">
						<template v-slot:activator="{ on }">
							<v-btn color="info" v-on="on">
								<v-icon>edit</v-icon>Change Password
							</v-btn>
						</template>

						<v-card>
							<v-card-title class="headline" primary-title>Change Password</v-card-title>

							<v-alert :value="error" type="warning" dismissible>{{error}}</v-alert>
							<v-alert :value="success" type="success" dismissible>{{success}}</v-alert>
							<v-card-text>
								<v-form>
									<v-text-field
										name="old_password"
										label="Old Password"
										type="password"
										v-model="changepass.old"
										single-line
									></v-text-field>
									<v-text-field
										name="new_password"
										label="New Password"
										type="password"
										v-model="changepass.new"
										single-line
									></v-text-field>
									<v-text-field
										name="repeat_new_password"
										label="Repeat New Password"
										type="password"
										v-model="changepass.rnew"
										single-line
									></v-text-field>
								</v-form>
							</v-card-text>

							<v-divider></v-divider>

							<v-card-actions>
								<v-btn @click="dialog=false">Cancel</v-btn>
								<v-spacer></v-spacer>
								<v-btn color="success" @click="ChangePassword">Change Password</v-btn>
							</v-card-actions>
						</v-card>
					</v-dialog>
				</v-layout>
			</v-flex>
		</v-layout>
	</v-container>
</template>


<script>
	export default {
		name: "Account",
		data: () => {
			return {
				dialog: false,
				changepass: {
					old: "",
					new: "",
					rnew: ""
				},
				error: "",
				success: ""
			};
		},
		methods: {
			LogOutAllDevices: function() {
				this.$http
					.post(
						this.$store.state.backendURL + "/api/v1/logout-all-devices/"
					)
					.then(resp => {
						if (resp.data.success) {
							this.$router.push("/login/");
						}
					});
			},
			ChangePassword: function() {
				this.$http
					.post(
						this.$store.state.backendURL + "/api/v1/change-password/",
						{
							old_password: this.changepass.old,
							new_password: this.changepass.new,
							repeat_new_password: this.changepass.rnew
						}
					)
					.then(resp => {
						if (resp.data.success) {
							this.success = "Sucesfully Changed the Password";
							this.changepass.old = "";
							this.changepass.new = "";
							this.changepass.rnew = "";
							setTimeout(this.CloseDialog, 1000);
						} else {
							this.error = resp.data.data.message;
						}
					})
					.catch(err => {
						this.error = err;
					});
			},
			CloseDialog: function() {
				this.dialog = false;
			}
		}
	};
</script>

<style scoped>
	.flex {
		padding-right: 2em;
	}
	.flex h1 {
		border-bottom: 1px solid white;
	}
	.v-icon {
		padding-right: 5px;
	}
</style>

