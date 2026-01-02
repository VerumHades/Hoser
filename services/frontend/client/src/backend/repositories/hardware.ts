import { backendRequest } from "../backend_request";

export interface ApiHardwareRate { 
    id: string
    resourceType: "CPU" | "RAM" | "DISK"
    costInCents: number
    validFromTime: string
}
export interface ApiHardwareRatesResponse { 
    rates: ApiHardwareRate[] 
}

export const HardwareAPI = {
    async getRates() { return backendRequest<ApiHardwareRatesResponse>("/hardware/rates", "GET"); },
};
