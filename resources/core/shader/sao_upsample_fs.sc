$input v_texcoord0

#include <bgfx_shader.sh>
#include <sao_common.sh>

SAMPLER2D(u_input, 1);
SAMPLER2D(u_depth_half, 2);

void main() {
    float a = 0.25;
    mat4 bilinear_weight = mat4(
        vec4(9,3,3,1),
        vec4(3,9,1,3),
        vec4(3,1,9,3),
        vec4(1,3,3,9)
    );

    float depth = linear_depth(texture2D(u_input, v_texcoord0).r);

    vec4 offset = floor((gl_FragCoord.xyxy + vec4( 1, 1,-1,-1)) / 2.0) * u_params[0].xyxy;

    vec4 source;
    source.x = texture2D(u_depth, offset.zw).r;
    source.y = texture2D(u_depth, offset.zy).r;
    source.z = texture2D(u_depth, offset.xw).r;
    source.w = texture2D(u_depth, offset.xy).r;

    vec4 half_res;
    half_res.x = texture2D(u_depth_half, offset.zw).r;
    half_res.y = texture2D(u_depth_half, offset.zy).r;
    half_res.z = texture2D(u_depth_half, offset.xw).r;
    half_res.w = texture2D(u_depth_half, offset.xy).r;

    int i = int(mod(gl_FragCoord.y, 2.0)) * 2 + int(mod(gl_FragCoord.x, 2.0));
    float w00 = bilinear_weight[i][0] / pow(1.0 + abs(half_res.x - depth), a);
    float w01 = bilinear_weight[i][1] / pow(1.0 + abs(half_res.y - depth), a);
    float w10 = bilinear_weight[i][2] / pow(1.0 + abs(half_res.z - depth), a);
    float w11 = bilinear_weight[i][3] / pow(1.0 + abs(half_res.w - depth), a);
            
    float w = w00 + w01 + w10 + w11;
    float ao = ((w00 * source.x) + (w01 * source.y) + (w10 * source.z) + (w11 * source.w)) / w;

    gl_FragColor = vec4(ao, ao, ao, 1.0);
}
