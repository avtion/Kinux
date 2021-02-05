// 判断时间戳是否有效
export function IsTimeOutLine(ttl: number): boolean {
  return Date.parse(new Date().toString()) / 1000 > ttl
}
