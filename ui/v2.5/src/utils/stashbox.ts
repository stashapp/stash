export function stashboxDisplayName(name: string, index: number) {
  return name || `Stash-Box #${index + 1}`;
}

export const getStashboxBase = (endpoint: string) =>
  endpoint.match(/(https?:\/\/.*?\/)graphql/)?.[1];
