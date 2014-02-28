module.exports = function(grunt) {

    grunt.initConfig({
        pkg: grunt.file.readJSON('package.json'),
        bower: grunt.file.readJSON('bower.json'),
        banner: '/**\n' +
                '* <%= bower.name %>.js v<%= bower.version %> \n' +
                '* <%= grunt.template.today("yyyy/mm/dd") %> \n' +
                '*/\n',
        shell: {
            goinstall: {
                options: {
                    failOnError: true,
                    stdout: true,
                    execOptions: {
                        cwd: '.'
                    }
                },
                command: 'go install -v ./...'
            }
        },
        concat: {
            options: {
                banner: '<%= banner %>',
                stripBanners: false
            },
            sitejs: {
                src: ['bower_components/jquery/dist/jquery.min.js', 'bower_components/angular/angular.min.js', 'bower_components/angular-animate/angular-animate.min.js', 'bower_components/angular-ui-bootstrap/index.js', 'bower_components/bootstrap/dist/js/bootstrap.min.js', 'js/<%= bower.name %>.js'],
                dest: 'dist/static/js/<%= bower.name %>.js'
            }
        },
        uglify: {
            options: {
                banner: '<%= banner %>'
            },
            sitejs: {
                files: {
                    'dist/static/js/<%= bower.name %>.min.js': ['<%= concat.sitejs.dest %>']
                }
            }
        },
        jshint: {
            options: {
                jshintrc: 'js/.jshintrc'
            },
            gruntfile: {
                src: 'Gruntfile.js'
            },
            src: {
                src: ['js/*.js']
            },
            test: {
                src: ['js/tests/unit/*.js']
            }
        },
        less: {
            compileCore: {
                options: {
                    strictMath: true,
                    sourceMap: true,
                    outputSourceFiles: true,
                    sourceMapURL: '<%= pkg.name %>.css.map',
                    sourceMapFilename: 'dist/static/css/<%= pkg.name %>.css.map'
                },
                files: {
                    'dist/static/css/<%= bower.name %>.css': ['less/<%= pkg.name %>.less']
                }
            },
            compileTheme: {
                options: {
                    strictMath: true,
                    sourceMap: true,
                    outputSourceFiles: true,
                    sourceMapURL: '<%= pkg.name %>-theme.css.map',
                    sourceMapFilename: 'dist/static/css/<%= pkg.name %>-theme.css.map'
                },
                files: {
                    'dist/static/css/<%= pkg.name %>-theme.css': 'less/theme.less'
                }
            },
            minify: {
                options: {
                    cleancss: true,
                    report: 'min'
                },
                files: {
                    'dist/static/css/<%= bower.name %>.min.css': 'dist/static/css/<%= bower.name %>.css'
                }
            }
        },
        copy: {
            images: {
                files: [
                    {
                        expand: true,
                        cwd: 'images/',
                        src: ['**/*.{png,jpg,gif}'],
                        dest: 'dist/static/images/'
                    }
                ]
            },
            templates: {
                files: [
                    {
                        src: ['templates/*', 'recipes', 'pages.json'],
                        dest: 'dist/'
                    }
                ]
            },
            bootstrap: {
                files: [
                    {
                        expand: true,
                        cwd: 'bower_components/bootstrap/dist/',
                        src: ['fonts/*'],
                        dest: 'dist/static/'
                    }
                ]
            },
            sfiles: {
                files: [
                    {
                        expand: true,
                        cwd: 'static/',
                        src: ['robots.txt', 'favicon.ico'],
                        dest: 'dist/static/'
                    }
                ]
            }
        }
    });

    require('load-grunt-tasks')(grunt, {scope: 'devDependencies'});

    grunt.registerTask('test', ['jshint']);
    grunt.registerTask('static-js', ['concat', 'uglify']);
    grunt.registerTask('static-css', ['less']);
    grunt.registerTask('static', ['copy', 'static-css', 'static-js']);

    grunt.registerTask('default', ['shell', 'test', 'static']);

};
