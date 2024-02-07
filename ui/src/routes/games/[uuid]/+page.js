export const load = ({ params }) => {
  console.log('loading page ', params)
  return { uuid: params.uuid }
}
