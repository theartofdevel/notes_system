import ax from "axios";
import config from "@/config";

export default ax.create({
  baseURL: config.API_URL,
});
