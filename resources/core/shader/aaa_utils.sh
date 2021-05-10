#ifndef AAA_UTILS_SH_HEADER_GUARD
#define AAA_UTILS_SH_HEADER_GUARD

#	if !defined(uv_ratio)
#		define uv_ratio vec2_splat(uAAAParams[0].x)
#	endif

vec2 NDCToViewRect(vec2 xy) { return ((xy * 0.5 + 0.5) * u_viewRect.zw + u_viewRect.xy) * uv_ratio; }

vec2 GetVelocityVector(in vec2 uv) {
#if BGFX_SHADER_LANGUAGE_GLSL
	const vec2 offset = vec2(0.5, 0.5);
#else
	const vec2 offset = vec2(0.5,-0.5);
#endif
	return texture2D(u_attr1, uv).xy * offset / (uResolution.xy / u_viewRect.zw);
}

#endif // AAA_UTILS_SH_HEADER_GUARD