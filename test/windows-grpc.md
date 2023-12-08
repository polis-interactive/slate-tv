grpcwebproxy --backend_addr=localhost:5000 --run_tls_server=false --use_websockets --allow_all_origins --server_http_debug_port=6969 --server_http_max_write_timeout=1h 

serve frontend -l 1234