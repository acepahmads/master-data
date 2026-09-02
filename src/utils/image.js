/**
 * Resolves full image URL from relative /uploads path or external web URL
 * @param {string} url 
 * @returns {string}
 */
export function getImageUrl(url) {
  if (!url) return '';
  if (url.startsWith('http://') || url.startsWith('https://') || url.startsWith('data:') || url.startsWith('blob:')) {
    return url;
  }
  if (url.startsWith('/uploads/')) {
    return `http://localhost:8080${url}`;
  }
  if (url.startsWith('uploads/')) {
    return `http://localhost:8080/${url}`;
  }
  return url;
}

export default getImageUrl;
