import axios from "@/utils/axios";

export interface ISignUpModel {
  email: string;
  password: string;
  repeatPassword: string;
}

export interface IAuthModel {
  email: string;
  password: string;
}

export const signUp = async (model: ISignUpModel) => {
  // try {
  //   const response = await axios.post("/api/signup", { ...model });
  //   console.log(response.status, response.data);
  // } catch (e) {}
};

export const auth = async (model: IAuthModel) => {
  // try {
  //   const response = await axios.post("/api/auth", { ...model });
  //   console.log(response.status, response.data);
  // } catch (e) {}
};
